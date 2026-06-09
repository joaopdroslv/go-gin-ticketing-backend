package user

import (
	"context"
	"go-gin-ticketing-backend/internal/domain"
	"go-gin-ticketing-backend/internal/shared/enums"
	"go-gin-ticketing-backend/internal/shared/schemas"
	"time"
)

type UserService struct {
	userRepository  UserRepository
	userStatusesMap map[int64]string
}

func NewUserService(ctx context.Context, userRepository UserRepository) (*UserService, error) {

	userStatuses, err := userRepository.GetAllUserStatuses(ctx)
	if err != nil {
		return nil, err
	}

	userStatusesMap := getUserStatusesMap(userStatuses)

	return &UserService{
		userRepository:  userRepository,
		userStatusesMap: userStatusesMap,
	}, nil
}

func getUserStatusesMap(userStatuses []UserStatus) map[int64]string {

	mapping := make(map[int64]string, len(userStatuses))

	for _, st := range userStatuses {
		mapping[st.ID] = st.Name
	}

	return mapping
}

func (s *UserService) GetAllUsers(
	ctx context.Context,
	paginationQuery schemas.PaginationQuery,
) (*GetAllUsersResponse, error) {

	pagination := domain.NewPagination(paginationQuery.Page, paginationQuery.Limit)

	users, total, err := s.userRepository.GetAllUsers(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return &GetAllUsersResponse{
		Items: s.translateUsers(users, s.userStatusesMap),
		Pagination: schemas.ResponsePagination{
			Page:      pagination.Page,
			PageTotal: int64(len(users)),
			Limit:     pagination.Limit,
			Total:     *total,
		},
	}, nil
}

func (s *UserService) GetAllUserStatuses(ctx context.Context) ([]UserStatus, error) {

	return s.userRepository.GetAllUserStatuses(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*ResponseUser, error) {

	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.translateUser(user, s.userStatusesMap), nil
}

func (s *UserService) CreateUser(ctx context.Context, body CreateUserBody) (*ResponseUser, error) {

	birthdate, err := time.Parse("2006-01-02", body.Birthdate)
	if err != nil {
		return nil, err
	}

	data := &CreateUserData{
		UserStatusID: int64(enums.PasswordCreationPending),
		Name:         body.Name,
		Birthdate:    birthdate,
		Email:        body.Email,
	}

	id, err := s.userRepository.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetUserByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return s.translateUser(user, s.userStatusesMap), nil
}

func (s *UserService) UpdateUserByID(
	ctx context.Context,
	id int64,
	body UpdateUserBody,
) (*ResponseUser, error) {

	data := &UpdateUserData{}

	if body.Name != nil {
		data.Name = body.Name
	}
	if body.Birthdate != nil {
		birthdate, err := time.Parse("2006-01-02", *body.Birthdate)
		if err != nil {
			return nil, err
		}
		data.Birthdate = &birthdate
	}
	if body.Email != nil {
		data.Email = body.Email
	}

	user, err := s.userRepository.UpdateUserByID(ctx, id, data)
	if err != nil {
		return nil, err
	}

	return s.translateUser(user, s.userStatusesMap), nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int64) (*DeleteUserResponse, error) {

	success, err := s.userRepository.DeleteUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeleteUserResponse{
		ID:      id,
		Deleted: success,
	}, nil
}

func (s *UserService) transformUserModelIntoResponseUser(
	userModel *User,
	userStatusName string,
) *ResponseUser {

	return &ResponseUser{
		ID:         userModel.ID,
		Name:       userModel.Name,
		Birthdate:  userModel.Birthdate.Format(time.RFC3339),
		Email:      userModel.Email,
		UserStatus: userStatusName,
	}
}

func (s *UserService) translateUser(
	userModel *User,
	userStatusesMap map[int64]string,
) *ResponseUser {

	return s.transformUserModelIntoResponseUser(userModel, userStatusesMap[userModel.UserStatusID])
}

func (s *UserService) translateUsers(
	userModels []User,
	userStatusesMap map[int64]string,
) []ResponseUser {

	responseUsers := make([]ResponseUser, 0, len(userModels))

	for _, u := range userModels {
		responseUsers = append(
			responseUsers,
			*s.transformUserModelIntoResponseUser(&u, userStatusesMap[u.UserStatusID]),
		)
	}

	return responseUsers
}
