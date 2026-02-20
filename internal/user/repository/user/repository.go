package user

import (
	"context"

	shareddomain "go-gin-ticketing-backend/internal/shared/domain"
	"go-gin-ticketing-backend/internal/user/domain"
	"go-gin-ticketing-backend/internal/user/schemas"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, pagination *shareddomain.Pagination) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUserByID(ctx context.Context, id int64, data schemas.UpdateUserBody) (*domain.User, error)
	DeleteUserByID(ctx context.Context, id int64) (bool, error)
}
