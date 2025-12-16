package user

import (
	"context"

	"ticket-io/internal/user/domain"
)

type Repository interface {
	GetAll(ctx context.Context) (*[]domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
}
