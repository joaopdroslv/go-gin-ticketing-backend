package accesscontrol

import (
	"context"
	models "go-gin-ticketing-backend/internal/access_control/models"
	"go-gin-ticketing-backend/internal/domain"
)

type PermissionRepository interface {
	GetAllPermissions(
		ctx context.Context, name *string, pagination *domain.Pagination,
	) ([]models.Permission, *int64, error)
	GetPermissionsByRoleID(ctx context.Context, userID int64) ([]models.Permission, error)
	UserHasPermission(ctx context.Context, id int64, permission string) (bool, error)
}
