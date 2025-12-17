package service

import (
	"context"
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/repository"
)

type UserStatusService struct {
	userStatusRepository repository.UserStatusRepository
}

func NewUserStatusService(
	userStatusRepository repository.UserStatusRepository,
) *UserStatusService {
	return &UserStatusService{userStatusRepository: userStatusRepository}
}

func (s *UserStatusService) GetAll(ctx context.Context) ([]domain.UserStatus, error) {
	return s.userStatusRepository.GetAll(ctx)
}

func (s *UserStatusService) GetStatusMap(ctx context.Context) (map[int64]string, error) {
	user_statuses, err := s.userStatusRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	mapping := make(map[int64]string, len(user_statuses))

	for _, st := range user_statuses {
		mapping[st.ID] = st.Name
	}

	return mapping, nil
}
