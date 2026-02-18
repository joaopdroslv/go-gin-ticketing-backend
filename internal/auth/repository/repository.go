package repository

import (
	"context"
	"ticket-io/internal/auth/domain"
)

type UserAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserAuth, error)
	RegisterUser(ctx context.Context, user *domain.UserAuth) (*domain.UserAuth, error)
	GetUserPermissions(ctx context.Context, userID int64) ([]domain.Permission, error)
}
