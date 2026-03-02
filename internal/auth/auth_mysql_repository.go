package auth

import (
	"context"
	"database/sql"
	"errors"
	"go-gin-ticketing-backend/internal/domain"
)

type AuthMysqlRepository struct {
	db *sql.DB
}

func NewAuthMysqlRepository(db *sql.DB) *AuthMysqlRepository {

	return &AuthMysqlRepository{db: db}
}

func (r *AuthMysqlRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*UserCredential, error) {

	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			user_credentials.email,
			user_credentials.password_hash,
			users.id,
			users.user_status_id
		FROM main.user_credentials
		JOIN main.users ON users.user_credential_id = user_credentials.id
		WHERE user_credentials.email = ?
		`,
		email,
	)

	var userCredential UserCredential

	if err := row.Scan(
		&userCredential.Email,
		&userCredential.PasswordHash,
		&userCredential.UserInfo.ID,
		&userCredential.UserInfo.UserStatusID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &userCredential, nil
}

func (r *AuthMysqlRepository) RegisterUser(
	ctx context.Context,
	data *RegisterUserData,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO main.user_credentials (email, password_hash) VALUES (?, ?)`,
		data.Email,
		data.PasswordHash,
	)
	if err != nil {
		return err
	}

	userCredentialID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	res, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO main.users (
			user_credential_id,
			user_status_id,
			name,
			birthdate
		) VALUES (?, ?, ? ,?)
		`,
		userCredentialID,
		data.UserStatusID,
		data.Name,
		data.Birthdate,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
