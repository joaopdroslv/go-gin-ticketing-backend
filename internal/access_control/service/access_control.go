package acservice

import "context"

type AccessControl interface {
	UserHasPermission(ctx context.Context, userID int64, requiredPermission string) (bool, error)
}
