package service

import "context"

type AccessControl interface {
	ValidateUserPermission(ctx context.Context, userID int64, permission string) (bool, error)
}
