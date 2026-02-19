package userstatus

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"go-gin-ticketing-backend/internal/user/models"
)

type mysqlUserStatusRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlUserStatusRepository {

	return &mysqlUserStatusRepository{db: db}
}

func (r *mysqlUserStatusRepository) ListUserStatuses(ctx context.Context) ([]models.UserStatus, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM main.user_statuses
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list user statuses query: %w", err)
	}
	defer rows.Close()

	userStatuses := make([]models.UserStatus, 0)

	for rows.Next() {
		var userStatus models.UserStatus

		if err := rows.Scan(
			&userStatus.ID,
			&userStatus.Name,
			&userStatus.Description,
			&userStatus.CreatedAt,
			&userStatus.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list user statuses scan: %w", err)
		}
		userStatuses = append(userStatuses, userStatus)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list user statuses rows: %w", err)
	}

	return userStatuses, nil
}
