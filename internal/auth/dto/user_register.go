package dto

type UserRegisterBody struct {
	Name      string
	Birthdate string
	Email     string
	Password  string
}

type UserRegisterResponse struct{}
