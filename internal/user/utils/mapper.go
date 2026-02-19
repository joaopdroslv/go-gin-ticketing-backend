package utils

import (
	"go-gin-ticketing-backend/internal/user/domain"
	"go-gin-ticketing-backend/internal/user/schemas"
	"time"
)

func FormatDomainUserToResponseUser(u *domain.User, userStatusName string) *schemas.ResponseUser {

	return &schemas.ResponseUser{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Birthdate:  u.Birthdate.Format(time.RFC3339),
		UserStatus: userStatusName,
	}
}

func DomainUserToResponseUser(u *domain.User, userStatusesMap map[int64]string) *schemas.ResponseUser {

	return FormatDomainUserToResponseUser(u, userStatusesMap[u.UserStatusID])
}

func DomainUsersToResponseUsers(domainUsers []domain.User, userStatusesMap map[int64]string) []schemas.ResponseUser {

	responseUsers := make([]schemas.ResponseUser, 0, len(domainUsers))

	for _, u := range domainUsers {
		responseUsers = append(responseUsers, *FormatDomainUserToResponseUser(&u, userStatusesMap[u.UserStatusID]))
	}

	return responseUsers
}
