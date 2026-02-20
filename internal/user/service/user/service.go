package user

import (
	"context"
	shareddomain "go-gin-ticketing-backend/internal/shared/domain"
	sharedschemas "go-gin-ticketing-backend/internal/shared/schemas"
	"go-gin-ticketing-backend/internal/user/domain"
	userrepository "go-gin-ticketing-backend/internal/user/repository/user"
	"go-gin-ticketing-backend/internal/user/schemas"
	"go-gin-ticketing-backend/internal/user/utils"
	"time"
)

type UserService struct {
	userRepository     userrepository.UserRepository
	userStatusProvider UserStatusProvider
}

func New(
	userRepository userrepository.UserRepository, userStatusProvider UserStatusProvider,
) *UserService {

	return &UserService{
		userRepository:     userRepository,
		userStatusProvider: userStatusProvider,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context, paginationQuery sharedschemas.PaginationQuery) (*schemas.GetAllUsersResponse, error) {

	pagination := shareddomain.NewPagination(paginationQuery.Page, paginationQuery.Limit)

	users, err := s.userRepository.GetAllUsers(ctx, pagination)
	if err != nil {
		return nil, err
	}

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	responseUsers := utils.DomainUsersToResponseUsers(users, userStatusesMap)

	return &schemas.GetAllUsersResponse{
		Total: int64(len(responseUsers)),
		Users: responseUsers,
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

	return utils.DomainUserToResponseUser(user, userStatusesMap), nil
}

func (s *UserService) CreateUser(ctx context.Context, body schemas.CreateUserBody) (*schemas.ResponseUser, error) {

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

	return utils.DomainUserToResponseUser(user, userStatusesMap), nil
}

func (s *UserService) UpdateUserByID(ctx context.Context, id int64, data schemas.UpdateUserBody) (*schemas.ResponseUser, error) {

	user, err := s.userRepository.UpdateUserByID(ctx, id, data)
	if err != nil {
		return nil, err
	}

	userStatusesMap, err := s.userStatusProvider.GetUserStatusesMap(ctx)
	if err != nil {
		return nil, err
	}

	return utils.DomainUserToResponseUser(user, userStatusesMap), nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int64) (*schemas.DeleteUserResponse, error) {

	success, err := s.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &schemas.DeleteUserResponse{
		ID:      id,
		Deleted: success,
	}, nil
}
