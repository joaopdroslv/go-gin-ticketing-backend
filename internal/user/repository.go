package user

import (
	"context"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) (*User, error)
}
