package services

import (
	"fmt"
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
)

func CreateInvoice(clientID int, amount float64, dueDate time.Time) (int, error) {

	// 1. cria invoice SEM payment_url
	var invoiceID int
	err := database.DB.QueryRow(`
		INSERT INTO invoices (client_id, amount, due_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, clientID, amount, dueDate).Scan(&invoiceID)

	if err != nil {
		return 0, err
	}

	// 2. cria checkout no Stripe usando invoiceID
	paymentURL, err := CreateCheckoutSession(amount, invoiceID)
	if err != nil {
		return 0, err
	}

	// 3. salva payment_url
	_, err = database.DB.Exec(`
		UPDATE invoices SET payment_url = $1 WHERE id = $2
	`, paymentURL, invoiceID)

	if err != nil {
		return 0, err
	}

	// 4. cria jobs automáticos
	createJobs(invoiceID)

	return invoiceID, nil
}

// 🔥 CRIA JOBS AUTOMÁTICOS
func createJobs(invoiceID int) {
	now := time.Now()

	jobs := []time.Time{
		now,
		now.Add(3 * 24 * time.Hour),
		now.Add(7 * 24 * time.Hour),
	}

	for _, runAt := range jobs {
		_, err := database.DB.Exec(`
			INSERT INTO jobs (invoice_id, run_at, type)
			VALUES ($1, $2, 'email')
		`, invoiceID, runAt)

		if err != nil {
			// em produção logaria
			fmt.Println("")
		}
	}
}
func MarkInvoiceAsPaid(invoiceID int) error {

	// 1. atualiza invoice
	_, err := database.DB.Exec(`
		UPDATE invoices
		SET status = 'paid'
		WHERE id = $1
	`, invoiceID)

	if err != nil {
		return err
	}

	// 2. cancela jobs futuros
	_, err = database.DB.Exec(`
		UPDATE jobs
		SET status = 'done'
		WHERE invoice_id = $1 AND status = 'pending'
	`, invoiceID)

	return err
}