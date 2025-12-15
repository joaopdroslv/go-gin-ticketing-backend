package user

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) (*[]User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}
