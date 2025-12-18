package status

import (
	"context"
	"ticket-io/internal/user/domain"
	statusrepository "ticket-io/internal/user/repository/status"
)

type StatusService struct {
	statusRepository statusrepository.StatusRepository
}

func New(statusRepository statusrepository.StatusRepository) *StatusService {

	return &StatusService{statusRepository: statusRepository}
}

func (s *StatusService) GetAll(ctx context.Context) ([]domain.Status, error) {

	return s.statusRepository.GetAll(ctx)
}

func (s *StatusService) GetStatusMap(ctx context.Context) (map[int64]string, error) {

	user_statuses, err := s.statusRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	mapping := make(map[int64]string, len(user_statuses))

	for _, st := range user_statuses {
		mapping[st.ID] = st.Name
	}

	return mapping, nil
}
