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
		SELECT j.id, j.invoice_id, j.type, c.email
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

		err := rows.Scan(&id, &invoiceID, &jobType, &email)
		if err != nil {
			log.Println("Error scanning job:", err)
			continue
		}

		processJob(id, invoiceID, jobType, email)
	}
}

func processJob(id int, invoiceID int, jobType string, email string) {
	log.Printf("Processing job %d for invoice %d (%s)", id, invoiceID, jobType)
	log.Println("Sending email to:", email)

	err := services.SendEmail(
		email,
		"Payment Reminder",
		"<h1>You have a pending invoice</h1>",
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