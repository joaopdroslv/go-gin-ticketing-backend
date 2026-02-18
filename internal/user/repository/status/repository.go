package status

import (
	"context"
	"ticket-io/internal/user/models"
)

type StatusRepository interface {
	ListStatuses(ctx context.Context) ([]models.Status, error)
}
