package permission

import (
	"context"
	"database/sql"
	"go-gin-ticketing-backend/internal/auth/models"
)

type mysqlPermissionRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *mysqlPermissionRepository {

	return &mysqlPermissionRepository{db: db}
}

func (r *mysqlPermissionRepository) GetPermissionsByUserID(ctx context.Context, id int64) ([]models.Permission, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			permissions.id,
			permissions.name,
			permissions.created_at,
			permissions.updated_at
		FROM permissions
		JOIN role_permissions ON role_permissions.permission_id = permissions.id
		JOIN user_roles ON user_roles.role_id = role_permissions.role_id
		JOIN users ON users.id = user_roles.user_id
		WHERE users.id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission

	for rows.Next() {
		var permission models.Permission

		if err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
