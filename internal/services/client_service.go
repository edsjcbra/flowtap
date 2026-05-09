package services

import "github.com/edsjcbra/flowtap/internal/database"

func CreateClient(
	name string,
	email string,
	phone string,
) error {

	_, err := database.DB.Exec(`
		INSERT INTO clients (
			name,
			email,
			phone
		)
		VALUES ($1, $2, $3)
	`,
		name,
		email,
		phone,
	)

	return err
}

func UpdateClient(
	id int,
	name string,
	email string,
	userID int,
) error {

	_, err := database.DB.Exec(`
		UPDATE clients
		SET
			name = $1,
			email = $2
		WHERE
			id = $3
		AND
			user_id = $4
	`,
		name,
		email,
		id,
		userID,
	)

	return err
}

func DeleteClient(id int) error {

	_, err := database.DB.Exec(`
		DELETE FROM clients
		WHERE id = $1
	`, id)

	return err
}