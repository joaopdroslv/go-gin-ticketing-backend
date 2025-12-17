package repository

import (
	"context"
	"ticket-io/internal/user/domain"
)

type UserStatusRepository interface {
	GetAll(ctx context.Context) ([]domain.UserStatus, error)
}
