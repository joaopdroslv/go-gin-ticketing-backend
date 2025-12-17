package repository

import (
	"context"

	"ticket-io/internal/user/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) (*[]domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
}
