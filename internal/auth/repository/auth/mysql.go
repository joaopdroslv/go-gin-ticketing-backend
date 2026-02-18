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

	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash FROM users WHERE email = ?`, email)

	var u domain.UserAuth

	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mysqlUserAuthRepository) RegisterUser(ctx context.Context, user *domain.UserAuth) (*domain.UserAuth, error) {

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, birthdate, email, password_hash, status_id) VALUES (?, ?, ?, ?, ?)`,
		user.Name, user.Birthdate, user.Email, user.PasswordHash, user.StatusID,
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	user.ID = id

	return user, nil
}
