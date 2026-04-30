// Package scheduler contains the rules to schedule the jobs invoices
package scheduler

import (
	"log"
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/services"
)

func Start() {
	go func() {
		for {
			runPendingJobs()
			time.Sleep(10 * time.Second)
		}
	}()
}

func runPendingJobs() {
	query := `
		SELECT j.id, j.invoice_id, j.type, c.email, i.payment_url
		FROM jobs j
		JOIN invoices i ON j.invoice_id = i.id
		JOIN clients c ON i.client_id = c.id
		WHERE j.status = 'pending'
		AND j.run_at <= NOW()
		AND i.status != 'paid'
	`

	rows, err := database.DB.Query(query)
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

	if err != nil {
		log.Println("EMAIL ERROR:", err)
		return
	}

	log.Println("EMAIL SENT SUCCESS")

	_, err = database.DB.Exec(`
		UPDATE jobs SET status = 'done' WHERE id = $1
	`, id)

	if err != nil {
		log.Println("Error updating job:", err)
	}
}
