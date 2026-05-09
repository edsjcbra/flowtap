package services

import (
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
)

// ================= CREATE INVOICE =================

func CreateInvoice(
	clientID int,
	amount float64,
	dueDate time.Time,
	userID int,
) (int, error) {

	var invoiceID int

	err := database.DB.QueryRow(`
		INSERT INTO invoices (
			client_id,
			amount,
			status,
			due_date,
			user_id
		)
		VALUES ($1, $2, 'pending', $3, $4)
		RETURNING id
	`,
		clientID,
		amount,
		dueDate,
		userID,
	).Scan(&invoiceID)

	if err != nil {
		return 0, err
	}

	now := time.Now()

	// FIRST EMAIL

	firstRun := dueDate

	if !dueDate.After(now) {
		firstRun = now.Add(2 * time.Second)
	}

	_, err = database.DB.Exec(`
		INSERT INTO jobs (
			invoice_id,
			run_at,
			type,
			status
		)
		VALUES ($1, $2, 'email', 'pending')
	`,
		invoiceID,
		firstRun,
	)

	if err != nil {
		return 0, err
	}

	// BASE DATE

	baseDate := dueDate

	if !dueDate.After(now) {
		baseDate = now
	}

	// +2 DAYS

	_, err = database.DB.Exec(`
		INSERT INTO jobs (
			invoice_id,
			run_at,
			type,
			status
		)
		VALUES ($1, $2, 'email', 'pending')
	`,
		invoiceID,
		baseDate.Add(2*24*time.Hour),
	)

	if err != nil {
		return 0, err
	}

	// +5 DAYS

	_, err = database.DB.Exec(`
		INSERT INTO jobs (
			invoice_id,
			run_at,
			type,
			status
		)
		VALUES ($1, $2, 'email', 'pending')
	`,
		invoiceID,
		baseDate.Add(5*24*time.Hour),
	)

	if err != nil {
		return 0, err
	}

	return invoiceID, nil
}

// ================= MARK AS PAID =================

func MarkInvoiceAsPaid(invoiceID int) error {

	// invoice -> paid

	_, err := database.DB.Exec(`
		UPDATE invoices
		SET status = 'paid'
		WHERE id = $1
	`, invoiceID)

	if err != nil {
		return err
	}

	// jobs -> paid

	_, err = database.DB.Exec(`
		UPDATE jobs
		SET status = 'paid'
		WHERE invoice_id = $1
	`, invoiceID)

	return err
}

// ================= CANCEL INVOICE =================

func CancelInvoice(invoiceID int) error {

	// invoice -> cancelled

	_, err := database.DB.Exec(`
		UPDATE invoices
		SET status = 'cancelled'
		WHERE id = $1
	`, invoiceID)

	if err != nil {
		return err
	}

	// jobs -> cancelled

	_, err = database.DB.Exec(`
		UPDATE jobs
		SET status = 'cancelled'
		WHERE invoice_id = $1
	`, invoiceID)

	return err
}

// ================= DELETE INVOICE =================

func DeleteInvoice(invoiceID int) error {

	// delete jobs first

	_, err := database.DB.Exec(`
		DELETE FROM jobs
		WHERE invoice_id = $1
	`, invoiceID)

	if err != nil {
		return err
	}

	// delete invoice

	_, err = database.DB.Exec(`
		DELETE FROM invoices
		WHERE id = $1
	`, invoiceID)

	return err
}