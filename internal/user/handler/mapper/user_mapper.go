package mapper

import (
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/dto"
	"time"
)

func UserToResponse(u *domain.User, statusName string) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Birthdate: u.Birthdate.Format(time.RFC3339),
		Status:    statusName,
	}
}

func UsersToResponse(users []domain.User, statusMap map[int64]string) []dto.UserResponse {
	res := make([]dto.UserResponse, 0, len(users))

	for _, u := range users {
		res = append(res, UserToResponse(&u, statusMap[u.StatusID]))
	}

	return res
}
