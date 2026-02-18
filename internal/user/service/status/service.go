package status

import (
	"context"
	"ticket-io/internal/user/models"
	statusrepository "ticket-io/internal/user/repository/status"
)

type StatusService struct {
	statusRepository statusrepository.StatusRepository
}

func New(statusRepository statusrepository.StatusRepository) *StatusService {

	return &StatusService{statusRepository: statusRepository}
}

func (s *StatusService) ListStatuses(ctx context.Context) ([]models.Status, error) {

	return s.statusRepository.ListStatuses(ctx)
}

func (s *StatusService) GetStatusMap(ctx context.Context) (map[int64]string, error) {

	userStatuses, err := s.statusRepository.ListStatuses(ctx)
	if err != nil {
		return nil, err
	}

	mapping := make(map[int64]string, len(userStatuses))

	for _, st := range userStatuses {
		mapping[st.ID] = st.Name
	}

	return mapping, nil
}
