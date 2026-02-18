package repository

import (
	"context"
	"ticket-io/internal/auth/domain"
	"ticket-io/internal/auth/models"
)

type UserAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserAuth, error)
	RegisterUser(ctx context.Context, user *domain.UserAuth) (*domain.UserAuth, error)
	GetUserPermissions(ctx context.Context, userID int64) ([]models.Permission, error)
}
