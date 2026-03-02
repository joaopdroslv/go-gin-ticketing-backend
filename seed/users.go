package seed

import (
	"database/sql"
	"go-gin-ticketing-backend/internal/shared/enums"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"
)

func Users(db *sql.DB, amount int) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userCredentialQuery := "INSERT INTO main.user_credentials (email, password_hash) VALUES (?, ?)"
	userQuery := `
		INSERT INTO main.users (
			user_credential_id,
			user_status_id,
			name,
			birthdate
		) VALUES (?, ?, ?, ?)
	`

	userCredentialStmt, err := tx.Prepare(userCredentialQuery)
	if err != nil {
		return err
	}
	defer userCredentialStmt.Close()

	userStmt, err := tx.Prepare(userQuery)
	if err != nil {
		return err
	}
	defer userStmt.Close()

	for i := 0; i < amount; i++ {

		birthdate := gofakeit.DateRange(
			time.Now().AddDate(-80, 0, 0),
			time.Now().AddDate(-18, 0, 0),
		)
		email := gofakeit.Email()
		password := gofakeit.Password(true, true, true, true, false, 12)
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

		res, err := userCredentialStmt.Exec(email, passwordHash)
		if err != nil {
			return err
		}

		userCredentialID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		_, err = userStmt.Exec(
			userCredentialID,
			int64(enums.Active),
			gofakeit.Name(),
			birthdate.Format("2006-01-02"),
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
