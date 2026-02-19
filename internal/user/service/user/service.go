package user

import (
	"context"
	"go-gin-ticketing-backend/internal/user/domain"
	"go-gin-ticketing-backend/internal/user/handler/mapper"
	userrepository "go-gin-ticketing-backend/internal/user/repository/user"
	"go-gin-ticketing-backend/internal/user/schemas"
	"time"
)

type UserService struct {
	userRepository     userrepository.UserRepository
	userStatusProvider UserStatusProvider
}

func New(
	userRepository userrepository.UserRepository, statusProvider UserStatusProvider,
) *UserService {

	return &UserService{
		userRepository:     userRepository,
		userStatusProvider: statusProvider,
	}
}

func (s *UserService) ListUsers(ctx context.Context) (*schemas.GetAllResponse, error) {

	users, err := s.userRepository.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	formattedUsers := mapper.UsersToResponseUsers(users, userStatusesMap)

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

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, userStatusesMap), nil
}

func (s *UserService) CreateUser(ctx context.Context, body schemas.UserCreateBody) (*schemas.ResponseUser, error) {

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(body.UserStatusID, body.Email, body.Name, birthdate)
	if err != nil {
		return nil, err
	}

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, userStatusesMap), nil
}

func (s *UserService) UpdateUserByID(ctx context.Context, id int64, data schemas.UserUpdateBody) (*schemas.ResponseUser, error) {

	user, err := s.userRepository.UpdateUserByID(ctx, id, data)
	if err != nil {
		return nil, err
	}

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.UserToResponseUser(user, userStatusesMap), nil
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
