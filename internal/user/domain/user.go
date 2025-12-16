package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, errors.New("e-mail is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &User{
		Email: email,
		Name:  name,
	}, nil
}
