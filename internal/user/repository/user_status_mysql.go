package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"ticket-io/internal/user/domain"
)

type mysqlUserStatusRepository struct {
	db *sql.DB
}

func NewMySQLUserStatusRepository(db *sql.DB) *mysqlUserStatusRepository {
	return &mysqlUserStatusRepository{db: db}
}

func (r *mysqlUserStatusRepository) GetAll(ctx context.Context) ([]domain.UserStatus, error) {

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM main.user_statuses ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user_statuses := make([]domain.UserStatus, 0)

	for rows.Next() {
		var s domain.UserStatus

		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		user_statuses = append(user_statuses, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user_statuses, nil
}
