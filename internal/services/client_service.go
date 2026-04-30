// Package services contain the business rules of the app
package services

import "github.com/edsjcbra/flowtap/internal/database"

func CreateClient(name, email, phone string) error {
	query := `
		INSERT INTO clients (name, email, phone)
		VALUES ($1, $2, $3)
	`

	_, err := database.DB.Exec(query, name, email, phone)
	return err
}
