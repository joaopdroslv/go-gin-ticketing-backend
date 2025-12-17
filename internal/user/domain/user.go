package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Name      string
	Birthdate time.Time
	StatusID  int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, name string, birthdate time.Time, statusID int64) (*User, error) {

	if email == "" {
		return nil, errors.New("e-mail is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if birthdate.IsZero() {
		return nil, errors.New("birthdate is required")
	}
	if statusID <= 0 {
		return nil, errors.New("status_id is required")
	}

	return &User{
		Email:     email,
		Name:      name,
		Birthdate: birthdate,
		StatusID:  statusID,
	}, nil
}
