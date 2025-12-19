package repository

import (
	"context"
	"database/sql"
	"ticket-io/internal/auth/domain"
)

type mysqlUserAuthRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlUserAuthRepository {

	return &mysqlUserAuthRepository{db: db}
}

func (r *mysqlUserAuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserAuth, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash FROM users WHERE email = ?
	`, email)

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

func (r *mysqlUserAuthRepository) GetUserPermissions(ctx context.Context, userID int64) (map[string]struct{}, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT permissions.name
		FROM permissions
		JOIN role_permissions ON role_permissions.permission_id = permissions.id
		JOIN user_roles ON user_roles.role_id = role_permissions.role_id
		JOIN users ON users.id = user_roles.user_id
		WHERE users.id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make(map[string]struct{})

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		permissions[name] = struct{}{}
	}

	return permissions, nil
}
