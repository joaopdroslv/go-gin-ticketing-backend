package mapper

import (
	"go-gin-ticketing-backend/internal/user/domain"
	"go-gin-ticketing-backend/internal/user/schemas"
	"time"
)

func FormatUserToResponseUser(u *domain.User, statusName string) *schemas.ResponseUser {

	return &schemas.ResponseUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Birthdate: u.Birthdate.Format(time.RFC3339),
		Status:    statusName,
	}
}

func UserToResponseUser(u *domain.User, statusMap map[int64]string) *schemas.ResponseUser {

	return FormatUserToResponseUser(u, statusMap[u.StatusID])
}

func UsersToResponseUsers(users []domain.User, statusMap map[int64]string) []schemas.ResponseUser {

	formattedUsers := make([]schemas.ResponseUser, 0, len(users))

	for _, u := range users {
		formattedUsers = append(formattedUsers, *FormatUserToResponseUser(&u, statusMap[u.StatusID]))
	}

	return formattedUsers
}
