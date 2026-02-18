package status

import (
	"context"
	"go-gin-ticketing-backend/internal/user/models"
)

type StatusRepository interface {
	ListStatuses(ctx context.Context) ([]models.Status, error)
}
