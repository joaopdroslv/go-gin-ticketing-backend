package service

import (
	"context"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/dto"
	"ticket-io/internal/user/repository"
	"time"
)

type UserService struct {
	userRepository    repository.UserRepository
	userStatusService UserStatusService
}

func NewUserService(
	userRepository repository.UserRepository,
	userStatusService *UserStatusService,
) *UserService {
	return &UserService{
		userRepository:    userRepository,
		userStatusService: *userStatusService,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {

	return s.userRepository.GetAll(ctx)
}

func (s *UserService) GetAllWithStatus(ctx context.Context) ([]domain.User, int64, map[int64]string, error) {

	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, 0, nil, err
	}

	statusMap, err := s.userStatusService.GetStatusMap(ctx)
	if err != nil {
		return nil, 0, nil, err
	}

	return users, int64(len(users)), statusMap, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {

	return s.userRepository.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, email, name string, birthdate time.Time, statusID int64) (*domain.User, error) {

	user, err := domain.NewUser(email, name, birthdate, statusID)
	if err != nil {
		return nil, err
	}

	return s.userRepository.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id int64, data dto.UserUpdateBody) (*domain.User, error) {
	return s.userRepository.Update(ctx, id, data)
}
