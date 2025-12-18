package user

import (
	"context"

	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/dto"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateByID(ctx context.Context, id int64, data dto.UserUpdateBody) (*domain.User, error)
	DeleteByID(ctx context.Context, id int64) (bool, error)
}
