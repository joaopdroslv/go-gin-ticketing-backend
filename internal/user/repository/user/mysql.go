package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"go-gin-ticketing-backend/internal/shared/errs"
	"go-gin-ticketing-backend/internal/user/domain"
	"go-gin-ticketing-backend/internal/user/schemas"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlUserRepository {

	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) ListUsers(ctx context.Context) ([]domain.User, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			users.id,
			users.user_status_id,
			users.email,
			users.name,
			users.birthdate,
			users.created_at,
			users.updated_at
		FROM main.users
		ORDER BY users.id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list users query: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var user domain.User

		if err := rows.Scan(
			&user.ID,
			&user.UserStatusID,
			&user.Email,
			&user.Name,
			&user.Birthdate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list users scan: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users rows error: %w", err)
	}

	return users, nil
}

func (r *mysqlUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {

	row := r.db.QueryRowContext(ctx, `
		SELECT
			users.id,
			users.user_status_id,
			users.email,
			users.name,
			users.birthdate,
			users.created_at,
			users.updated_at
		FROM main.users
		WHERE users.id = ?
	`, id)

	var user domain.User

	if err := row.Scan(
		&user.ID,
		&user.UserStatusID,
		&user.Email,
		&user.Name,
		&user.Birthdate,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user id=%d: %w", id, errs.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("get user id=%d scan: %w", id, err)
	}

	return &user, nil
}

func (r *mysqlUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (
			users.user_status_id,
			users.email,
			users.name,
			users.birthdate
		) VALUES (?, ?, ?, ?)`,
		user.UserStatusID, user.Email, user.Name, user.Birthdate,
	)
	if err != nil {
		return nil, fmt.Errorf("create user exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create user last insert id: %w", err)
	}

	user.ID = id

	return user, nil
}

func (r *mysqlUserRepository) UpdateUserByID(ctx context.Context, id int64, data schemas.UserUpdateBody) (*domain.User, error) {

	fields := []string{}
	args := []any{}

	if data.Name != nil {
		fields = append(fields, "users.name = ?")
		args = append(args, data.Name)
	}

	if data.Email != nil {
		fields = append(fields, "users.email = ?")
		args = append(args, data.Email)
	}

	if data.Birthdate != nil {
		fields = append(fields, "users.birthdate = ?")
		args = append(args, data.Birthdate)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("update user: %w", errs.ErrNothingToUpdate)
	}

	query := fmt.Sprintf(
		"UPDATE main.users SET %s WHERE users.id = ?",
		strings.Join(fields, ", "),
	)

	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("update user id=%d exec: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update user id=%d rows affected: %w", id, err)
	}

	if rows == 0 {
		return nil, fmt.Errorf("user id=%d: %w", id, errs.ErrResourceNotFound)
	}

	return r.GetUserByID(ctx, id)
}

func (r *mysqlUserRepository) DeleteUserByID(ctx context.Context, id int64) (bool, error) {

	result, err := r.db.ExecContext(ctx, `DELETE FROM main.users WHERE users.id = ?`, id)
	if err != nil {
		return false, fmt.Errorf("delete user id=%d exec: %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete user id=%d rows affected: %w", id, err)
	}

	if rows == 0 {
		return false, fmt.Errorf("user id=%d: %w", id, errs.ErrResourceNotFound)
	}

	return true, nil
}
