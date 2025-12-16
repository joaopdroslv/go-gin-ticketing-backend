package user

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"ticket-io/internal/user/domain"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *mysqlRepository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetAll(ctx context.Context) (*[]domain.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, name, created_at, updated_at FROM users ORDER BY id DESC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *mysqlRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, created_at, updated_at FROM users WHERE id = ?`, id)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *mysqlRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (email, name) VALUES (?, ?) RETURNING id`,
		user.Email, user.Name,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}
