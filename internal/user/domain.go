package user

import "errors"

type User struct {
	ID    string
	Email string
	Name  string
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
