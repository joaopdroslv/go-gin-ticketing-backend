package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	shareddoamin "go-gin-ticketing-backend/internal/shared/domain"
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

func (r *mysqlUserRepository) GetAllUsers(ctx context.Context, pagination *shareddoamin.Pagination) ([]domain.User, error) {

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
		LIMIT ?
		OFFSET ?
	`, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func (r *mysqlUserRepository) UpdateUserByID(ctx context.Context, id int64, data schemas.UpdateUserBody) (*domain.User, error) {

	query, args, err := r.formatUpdateUserQuery(id, data)
	if err != nil {
		return nil, err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, errs.ErrZeroRowsAffected
	}

	return r.GetUserByID(ctx, id)
}

func (r mysqlUserRepository) formatUpdateUserQuery(id int64, data schemas.UpdateUserBody) (string, []any, error) {

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
		return "", nil, fmt.Errorf("update user: %w", errs.ErrNothingToUpdate)
	}

	query := fmt.Sprintf(
		"UPDATE main.users SET %s WHERE users.id = ?",
		strings.Join(fields, ", "),
	)

	args = append(args, id)

	return query, args, nil
}

func (r *mysqlUserRepository) DeleteUserByID(ctx context.Context, id int64) (bool, error) {

	result, err := r.db.ExecContext(ctx, `DELETE FROM main.users WHERE users.id = ?`, id)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, errs.ErrZeroRowsAffected
	}

	return true, nil
}
