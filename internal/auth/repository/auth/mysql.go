package auth

import (
	"context"
	"database/sql"
	"go-gin-ticketing-backend/internal/auth/domain"
)

type mysqlUserAuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlUserAuthRepository {

	return &mysqlUserAuthRepository{db: db}
}

func (r *mysqlUserAuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuth, error) {

	row := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			user_status_id,
			email,
			password_hash
		FROM users
		WHERE email = ?
	`, email)

	var userAuth domain.UserAuth

	if err := row.Scan(
		&userAuth.ID,
		&userAuth.UserStatusID,
		&userAuth.Email,
		&userAuth.PasswordHash,
	); err != nil {
		return nil, err
	}

	return &userAuth, nil
}

func (r *mysqlUserAuthRepository) RegisterUser(ctx context.Context, user *domain.UserAuth) (*domain.UserAuth, error) {

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (user_status_id, name, birthdate, email, password_hash) VALUES (?, ?, ?, ?, ?)`,
		user.UserStatusID, user.Name, user.Birthdate, user.Email, user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	user.ID = id

	return user, nil
}
