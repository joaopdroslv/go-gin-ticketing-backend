package accesscontrol

import (
	"context"
	"database/sql"
	models "go-gin-ticketing-backend/internal/access_control/models"
	"go-gin-ticketing-backend/internal/domain"
	"strings"
)

type PermissionRepositoryMysql struct {
	db *sql.DB
}

func NewPermissionRepositoryMysql(db *sql.DB) *PermissionRepositoryMysql {

	return &PermissionRepositoryMysql{db: db}
}

func (r *PermissionRepositoryMysql) GetAllPermissions(
	ctx context.Context,
	name *string,
	pagination *domain.Pagination,
) ([]models.Permission, *int64, error) {

	query, args := r.formatGetAllPermissionsQuery(name, pagination)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	var total int64

	for rows.Next() {
		var permission models.Permission
		var totalCount int64

		if err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&totalCount,
		); err != nil {
			return nil, nil, err
		}

		if total == 0 {
			total = totalCount
		}

		permissions = append(permissions, permission)
	}

	return permissions, &total, nil
}

func (r *PermissionRepositoryMysql) formatGetAllPermissionsQuery(
	name *string,
	pagination *domain.Pagination,
) (string, []any) {

	conditions := []string{}
	args := []any{}

	if name != nil && *name != "" {
		conditions = append(conditions, "permissions.name LIKE ?")
		args = append(args, "%"+*name+"%")
	}

	query := `
		SELECT
			permissions.id,
			permissions.name,
			permissions.description,
			permissions.created_at,
			permissions.updated_at,
			COUNT(*) OVER() AS total_count
		FROM main.permissions
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		ORDER BY permissions.id DESC
		LIMIT ?
		OFFSET ?
	`

	args = append(args, pagination.Limit)
	args = append(args, pagination.Offset)

	return query, args
}

func (r *PermissionRepositoryMysql) GetPermissionsByRoleID(
	ctx context.Context,
	id int64,
) ([]models.Permission, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			permissions.id,
			permissions.name,
			permissions.description,
			permissions.created_at,
			permissions.updated_at
		FROM main.permissions
		JOIN main.role_permissions ON role_permissions.permission_id = permissions.id
		JOIN main.roles ON roles.id = role_permissions.role_id
		WHERE roles.id = ?
		`,
		id,
	)
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
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *PermissionRepositoryMysql) UserHasPermission(
	ctx context.Context,
	id int64,
	permission string,
) (bool, error) {

	var exists int64

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT 1
		FROM main.users
		JOIN main.user_roles ON user_roles.user_id = users.id
		JOIN main.role_permissions ON role_permissions.role_id = user_roles.role_id
		JOIN main.permissions ON permissions.id = role_permissions.permission_id
		WHERE TRUE
			AND users.id = ?
			AND permissions = '?'
		LIMIT 1
		`,
		id, permission,
	).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
