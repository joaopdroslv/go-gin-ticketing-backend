package user

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *mysqlRepository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) FindByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, name FROM users WHERE id = ?`, id)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mysqlRepository) Save(ctx context.Context, user *User) (*User, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, name) VALUES (?, ?) RETURNING id`,
		user.Email, user.Name,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}
