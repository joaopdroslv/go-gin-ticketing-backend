package status

import (
	"context"
	"ticket-io/internal/user/domain"
)

type StatusRepository interface {
	GetAll(ctx context.Context) ([]domain.Status, error)
}
