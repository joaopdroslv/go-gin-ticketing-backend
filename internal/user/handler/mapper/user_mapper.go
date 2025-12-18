package mapper

import (
	"ticket-io/internal/user/domain"
	"ticket-io/internal/user/handler/dto"
	"time"
)

func UserToResponseUser(u *domain.User, statusName string) dto.ResponseUser {
	return dto.ResponseUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Birthdate: u.Birthdate.Format(time.RFC3339),
		Status:    statusName,
	}
}

func UsersToResponse(users []domain.User, statusMap map[int64]string) []dto.ResponseUser {
	res := make([]dto.ResponseUser, 0, len(users))

	for _, u := range users {
		res = append(res, UserToResponseUser(&u, statusMap[u.StatusID]))
	}

	return res
}
