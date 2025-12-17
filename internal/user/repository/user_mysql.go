package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"ticket-io/internal/user/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *mysqlUserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetAll(ctx context.Context) (*[]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			email,
			name,
			birthdate,
			status_id,
			created_at,
			updated_at
		FROM users
		ORDER BY id DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Name,
			&u.Birthdate,
			&u.StatusID,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			email,
			name,
			birthdate,
			status_id,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`, id)

	var u domain.User

	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Birthdate,
		&u.StatusID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, name, birthdate, status_id) VALUES (?, ?, ?, ?)`,
		user.Email,
		user.Name,
		user.Birthdate,
		user.StatusID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

// func (r *mysqlRepository) ChangeStatusByID(ctx context.Context, id int)
