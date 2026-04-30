// Package scheduler contains the rules to schedule the jobs invoices
package scheduler

import (
	"log"
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
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
		SELECT j.id, j.invoice_id, j.type
		FROM jobs j
		JOIN invoices i ON j.invoice_id = i.id
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

		err := rows.Scan(&id, &invoiceID, &jobType)
		if err != nil {
			log.Println("Error scanning job:", err)
			continue
		}

		processJob(id, invoiceID, jobType)
	}
}

func processJob(id int, invoiceID int, jobType string) {
	log.Printf("Processing job %d for invoice %d (%s)", id, invoiceID, jobType)

	// aqui depois entra email/SMS real

	_, err := database.DB.Exec(`
		UPDATE jobs SET status = 'done' WHERE id = $1
	`, id)

	if err != nil {
		log.Println("Error updating job:", err)
	}
}