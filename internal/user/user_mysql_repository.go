package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"go-gin-ticketing-backend/internal/domain"
)

type UserMysqlRepository struct {
	db *sql.DB
}

func NewUserMysqlRepository(db *sql.DB) *UserMysqlRepository {

	return &UserMysqlRepository{db: db}
}

func (r *UserMysqlRepository) GetAllUsers(
	ctx context.Context,
	pagination *domain.Pagination,
) ([]User, *int64, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			users.id,
			users.user_credential_id,
			users.user_status_id,
			users.name,
			users.birthdate,
			user_credentials.email,
			users.created_at,
			users.updated_at,
			COUNT(*) OVER() AS total_count
		FROM main.users
		JOIN main.user_credentials ON user_credentials.id = users.user_credential_id
		ORDER BY users.id DESC
		LIMIT ?
		OFFSET ?
		`,
		pagination.Limit,
		pagination.Offset,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var users []User
	var total int64

	for rows.Next() {
		var user User
		var totalCount int64

		if err := rows.Scan(
			&user.ID,
			&user.UserCredentialID,
			&user.UserStatusID,
			&user.Name,
			&user.Birthdate,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&totalCount,
		); err != nil {
			return nil, nil, err
		}

		if total == 0 {
			total = totalCount
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return users, &total, nil
}

func (r *UserMysqlRepository) GetAllUserStatuses(ctx context.Context) ([]UserStatus, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			user_statuses.id,
			user_statuses.name,
			user_statuses.description,
			user_statuses.created_at,
			user_statuses.updated_at
		FROM main.user_statuses
		ORDER BY user_statuses.id DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userStatuses []UserStatus

	for rows.Next() {
		var userStatus UserStatus

		if err := rows.Scan(
			&userStatus.ID,
			&userStatus.Name,
			&userStatus.Description,
			&userStatus.CreatedAt,
			&userStatus.UpdatedAt,
		); err != nil {
			return nil, err
		}

		userStatuses = append(userStatuses, userStatus)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userStatuses, nil
}

func (r *UserMysqlRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {

	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			users.id,
			users.user_credential_id,
			users.user_status_id,
			users.name,
			users.birthdate,
			user_credentials.email,
			users.created_at,
			users.updated_at
		FROM main.users
		JOIN main.user_credentials ON user_credentials.id = users.user_credential_id
		WHERE users.id = ?
		`,
		id,
	)

	var user User

	if err := row.Scan(
		&user.ID,
		&user.UserCredentialID,
		&user.UserStatusID,
		&user.Name,
		&user.Birthdate,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserMysqlRepository) CreateUser(
	ctx context.Context,
	data *CreateUserData,
) (*int64, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO main.user_credentials (email) VALUES (?)`,
		data.Email,
	)
	if err != nil {
		return nil, err
	}

	userCredentialID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	result, err = tx.ExecContext(
		ctx,
		`INSERT INTO main.users (
			user_credential_id,
			user_status_id,
			name,
			birthdate
		) VALUES (?, ?, ?, ?)`,
		userCredentialID,
		data.UserStatusID,
		data.Name,
		data.Birthdate,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *UserMysqlRepository) UpdateUserByID(
	ctx context.Context,
	id int64,
	data *UpdateUserData,
) (*User, error) {

	query, args, err := r.formatUpdateUserByIDQuery(id, data)
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
		return nil, domain.ErrUserNotFound
	}

	return r.GetUserByID(ctx, id)
}

func (r UserMysqlRepository) formatUpdateUserByIDQuery(
	id int64,
	data *UpdateUserData,
) (string, []any, error) {

	userFields := []string{}
	userCredentialFields := []string{}
	args := []any{}

	if data.Name != nil {
		userFields = append(userFields, "users.name = ?")
		args = append(args, data.Name)
	}
	if data.Birthdate != nil {
		userFields = append(userFields, "users.birthdate = ?")
		args = append(args, data.Birthdate)
	}
	if data.Email != nil {
		userCredentialFields = append(userCredentialFields, "user_credentials.email = ?")
		args = append(args, data.Email)
	}

	if len(userFields) == 0 && len(userCredentialFields) == 0 {
		return "", nil, domain.ErrNothingToUpdate
	}

	setItems := []string{}
	setItems = append(setItems, userFields...)
	setItems = append(setItems, userCredentialFields...)

	query := "UPDATE main.users"

	if len(userCredentialFields) > 0 {
		query += " JOIN main.user_credentials ON user_credentials.id = users.user_credential_id"
	}

	query += fmt.Sprintf(" SET %s WHERE users.id = ?", strings.Join(setItems, ", "))
	args = append(args, id)

	log.Println(query)

	return query, args, nil
}

func (r *UserMysqlRepository) DeleteUserByID(ctx context.Context, id int64) (bool, error) {

	result, err := r.db.ExecContext(ctx, `DELETE FROM main.users WHERE users.id = ?`, id)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, domain.ErrUserNotFound
	}

	return true, nil
}
