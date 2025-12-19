package domain

import (
	"errors"
	"time"
)

type UserAuth struct {
	ID           int64
	Name         string
	Birthdate    time.Time
	StatusID     int64
	Email        string
	PasswordHash string
}

func NewUserAuth(name string, birthdate time.Time, statusID int64, email, passwordHash string) (*UserAuth, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}
	if birthdate.IsZero() {
		return nil, errors.New("birthdate is required")
	}
	if statusID <= 0 {
		return nil, errors.New("statusID is required")
	}
	if email == "" {
		return nil, errors.New("e-mail is required")
	}
	if passwordHash == "" {
		return nil, errors.New("passwordHash is required")
	}

	return &UserAuth{
		Name:         name,
		Birthdate:    birthdate,
		StatusID:     statusID,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
