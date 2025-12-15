package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context) (*[]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, email, name string) (*User, error) {
	user, err := NewUser(email, name)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, user)
}
