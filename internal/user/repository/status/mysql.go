package status

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"ticket-io/internal/user/domain"
)

type mysqlStatusRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlStatusRepository {

	return &mysqlStatusRepository{db: db}
}

func (r *mysqlStatusRepository) GetAll(ctx context.Context) ([]domain.Status, error) {

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM main.user_statuses ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user_statuses := make([]domain.Status, 0)

	for rows.Next() {
		var s domain.Status

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
