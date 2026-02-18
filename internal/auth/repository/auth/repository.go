package auth

import (
	"context"
	"go-gin-ticketing-backend/internal/auth/domain"
)

type UserAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.UserAuth, error)
	RegisterUser(ctx context.Context, user *domain.UserAuth) (*domain.UserAuth, error)
}
