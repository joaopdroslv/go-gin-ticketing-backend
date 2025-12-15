package seed

import (
	"database/sql"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedUsersTable(db *sql.DB, amount int) error {
	query := `
		INSERT INTO users (name, email, birthdate)
		VALUES (?, ?, ?)
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
		); err != nil {
			return err
		}
	}

	return nil
}
