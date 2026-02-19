package userstatus

import (
	"context"
	"go-gin-ticketing-backend/internal/user/models"
	userstatusrepository "go-gin-ticketing-backend/internal/user/repository/user_status"
)

type UserStatusService struct {
	userStatusRepository userstatusrepository.UserStatusRepository
}

func New(userStatusRepository userstatusrepository.UserStatusRepository) *UserStatusService {

	return &UserStatusService{userStatusRepository: userStatusRepository}
}

func (s *UserStatusService) ListUserStatuses(ctx context.Context) ([]models.UserStatus, error) {

	return s.userStatusRepository.ListUserStatuses(ctx)
}

func (s *UserStatusService) GetUserStatusesMap(ctx context.Context) (map[int64]string, error) {

	userStatuses, err := s.userStatusRepository.ListUserStatuses(ctx)
	if err != nil {
		return nil, err
	}

	mapping := make(map[int64]string, len(userStatuses))

	for _, st := range userStatuses {
		mapping[st.ID] = st.Name
	}

	return mapping, nil
}
