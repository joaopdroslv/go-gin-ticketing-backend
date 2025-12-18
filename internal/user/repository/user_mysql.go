package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"ticket-io/internal/shared/errors"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/dto"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *mysqlUserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {

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

	return users, nil
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {

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

func (r *mysqlUserRepository) UpdateByID(ctx context.Context, id int64, data dto.UserUpdateBody) (*domain.User, error) {

	fields := []string{}
	args := []any{}

	if data.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *data.Name)
	}

	if data.Email != nil {
		fields = append(fields, "email = ?")
		args = append(args, *data.Email)
	}

	if data.Birthdate != nil {
		fields = append(fields, "birthdate = ?")
		args = append(args, *data.Birthdate)
	}

	if len(fields) == 0 {
		return nil, errors.ErrNothingToUpdate
	}

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = ?",
		strings.Join(fields, ", "),
	)

	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, errors.ErrZeroRowsAffected
	}

	return r.GetByID(ctx, id)
}

func (r *mysqlUserRepository) DeleteByID(ctx context.Context, id int64) (bool, error) {

	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, nil
	}

	if rows == 0 {
		return false, errors.ErrZeroRowsAffected
	}

	return true, nil
}

// func (r *mysqlUserRepository) ChangeStatusByID(ctx context.Context, id int64) (*domain.User, error) {}
