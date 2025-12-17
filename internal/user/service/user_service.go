package service

import (
	"context"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/repository"
	"time"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{userRepository: r}
}

func (s *UserService) GetAll(ctx context.Context) (*[]domain.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, email, name string, birthdate time.Time, statusID int64) (*domain.User, error) {
	user, err := domain.NewUser(email, name, birthdate, statusID)
	if err != nil {
		return nil, err
	}

	return s.userRepository.Create(ctx, user)
}
