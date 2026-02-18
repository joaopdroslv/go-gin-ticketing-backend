package permission

import (
	"context"
	"go-gin-ticketing-backend/internal/auth/models"
)

type PermissionRepository interface {
	GetPermissionsByUserID(ctx context.Context, userID int64) ([]models.Permission, error)
}
