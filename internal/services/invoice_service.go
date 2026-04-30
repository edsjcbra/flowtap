package services

import (
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
)

func CreateInvoice(clientID int, amount float64, dueDate time.Time) (int, error) {
	query := `
		INSERT INTO invoices (client_id, amount, due_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var invoiceID int
	err := database.DB.QueryRow(query, clientID, amount, dueDate).Scan(&invoiceID)
	if err != nil {
		return 0, err
	}


	createJobs(invoiceID)

	return invoiceID, nil
}

func createJobs(invoiceID int) {
	now := time.Now()

	jobs := []time.Time{
		now,
		now.Add(3 * 24 * time.Hour),
		now.Add(7 * 24 * time.Hour),
	}

	for _, runAt := range jobs {
		query := `
			INSERT INTO jobs (invoice_id, run_at, type)
			VALUES ($1, $2, 'email')
		`
		database.DB.Exec(query, invoiceID, runAt)
	}
}
