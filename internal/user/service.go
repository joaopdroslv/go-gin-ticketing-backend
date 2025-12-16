package user

import (
	"context"
	"ticket-io/internal/user/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context) (*[]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, email, name string) (*domain.User, error) {
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}
