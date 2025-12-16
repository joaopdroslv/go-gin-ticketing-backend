package service

import (
	"context"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/repository"
	"time"
)

type Service struct {
	repository repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) GetAll(ctx context.Context) (*[]domain.User, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, email, name string, birthdate time.Time, statusID int64) (*domain.User, error) {
	user, err := domain.NewUser(email, name, birthdate, statusID)
	if err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, user)
}
