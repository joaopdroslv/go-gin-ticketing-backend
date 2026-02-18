package user

import (
	"context"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/mapper"
	userrepository "ticket-io/internal/user/repository/user"
	"ticket-io/internal/user/schemas"
	"time"
)

type UserService struct {
	userRepository userrepository.UserRepository
	statusProvider StatusProvider
}

func New(
	userRepository userrepository.UserRepository, statusProvider StatusProvider,
) *UserService {

	return &UserService{
		userRepository: userRepository,
		statusProvider: statusProvider,
	}
}

func (s *UserService) ListUsers(ctx context.Context) (*schemas.GetAllResponse, error) {

	users, err := s.userRepository.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	statusMap, err := s.statusProvider.GetStatusMap(ctx)
	if err != nil {
		return nil, err
	}

	formattedUsers := mapper.UsersToResponseUsers(users, statusMap)

	return &schemas.GetAllResponse{
		Total: int64(len(formattedUsers)),
		Users: formattedUsers,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*schemas.ResponseUser, error) {

	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	statusMap, err := s.statusProvider.GetStatusMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, statusMap), nil
}

func (s *UserService) CreateUser(ctx context.Context, body schemas.UserCreateBody) (*schemas.ResponseUser, error) {

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(body.Email, body.Name, birthdate, body.StatusID)
	if err != nil {
		return nil, err
	}

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	statusMap, err := s.statusProvider.GetStatusMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, statusMap), nil
}

func (s *UserService) UpdateUserByID(ctx context.Context, id int64, data schemas.UserUpdateBody) (*schemas.ResponseUser, error) {

	user, err := s.userRepository.UpdateUserByID(ctx, id, data)
	if err != nil {
		return nil, err
	}

	statusMap, err := s.statusProvider.GetStatusMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, statusMap), nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int64) (*schemas.UserDeleteResponse, error) {

	success, err := s.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &schemas.UserDeleteResponse{
		ID:      id,
		Deleted: success,
	}, nil
}
