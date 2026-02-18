package status

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"ticket-io/internal/user/models"
)

type mysqlStatusRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlStatusRepository {

	return &mysqlStatusRepository{db: db}
}

func (r *mysqlStatusRepository) ListStatuses(ctx context.Context) ([]models.Status, error) {

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM main.user_statuses ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list user statuses query: %w", err)
	}
	defer rows.Close()

	userStatuses := make([]models.Status, 0)

	for rows.Next() {
		var s models.Status

		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list user statuses scan: %w", err)
		}
		userStatuses = append(userStatuses, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list user statuses rows error: %w", err)
	}

	return userStatuses, nil
}
