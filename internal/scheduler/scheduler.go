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
		SELECT id, invoice_id, type
		FROM jobs
		WHERE status = 'pending'
		AND run_at <= NOW()
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
			continue
		}

		processJob(id, invoiceID, jobType)
	}
}

func processJob(id int, invoiceID int, jobType string) {
	log.Printf("Processing job %d for invoice %d (%s)", id, invoiceID, jobType)

	// 🔥 por enquanto só simula envio
	// depois vamos integrar email/SMS

	_, err := database.DB.Exec(`
		UPDATE jobs SET status = 'done' WHERE id = $1
	`, id)

	if err != nil {
		log.Println("Error updating job:", err)
	}
}