// Package worker process jobs
package worker

import (
	"log"
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/services"
)

func StartWorker() {
	for {
		processPendingJobs()
		time.Sleep(5 * time.Second)
	}
}

func processPendingJobs() {

	rows, err := database.DB.Query(`
		SELECT j.id, j.invoice_id, j.type,
       	c.email,
       	COALESCE(i.payment_url, '')
		FROM jobs j
		JOIN invoices i ON j.invoice_id = i.id
		JOIN clients c ON i.client_id = c.id
		WHERE j.status = 'pending'
		AND j.run_at <= NOW()
		AND i.status != 'paid'
		LIMIT 10
	`)

	if err != nil {
		log.Println("Error fetching jobs:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var id int
		var invoiceID int
		var jobType string
		var email string
		var paymentURL string

		err := rows.Scan(&id, &invoiceID, &jobType, &email, &paymentURL)
		if err != nil {
			log.Println("Error scanning job:", err)
			continue
		}

		processJob(id, invoiceID, jobType, email, paymentURL)
	}
}

func processJob(id int, invoiceID int, jobType string, email string, paymentURL string) {
	log.Printf("Processing job %d for invoice %d (%s)", id, invoiceID, jobType)
	log.Println("Sending email to:", email)

	body := `
		<h2>Payment Reminder</h2>
		<p>You have a pending invoice.</p>
		<p><a href="` + paymentURL + `">👉 Pay Now</a></p>
	`

	err := services.SendEmail(
		email,
		"Invoice Payment Reminder",
		body,
	)

	// 🔥 SE DER ERRO → NÃO FINALIZA JOB
	if err != nil {
		log.Println("EMAIL ERROR:", err)
		return
	}

	log.Println("EMAIL SENT SUCCESS")

	// ✅ SÓ MARCA COMO DONE SE DEU CERTO
	_, err = database.DB.Exec(`
		UPDATE jobs SET status = 'done' WHERE id = $1
	`, id)

	if err != nil {
		log.Println("Error updating job:", err)
	}
}
