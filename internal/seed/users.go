package seed

import (
	"database/sql"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func Users(db *sql.DB, amount int) error {

	var statusID int64

	err := db.QueryRow(`
		SELECT id FROM user_statuses WHERE name = 'active'
	`).Scan(&statusID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, birthdate, status_id)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < amount; i++ {
		name := gofakeit.Name()
		email := gofakeit.Email()

		birthdate := gofakeit.DateRange(
			time.Now().AddDate(-80, 0, 0),
			time.Now().AddDate(-18, 0, 0),
		)

		if _, err := stmt.Exec(
			name,
			email,
			birthdate.Format("2006-01-02"),
			statusID,
		); err != nil {
			return err
		}
	}

	return nil
}
