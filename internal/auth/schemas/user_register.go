package schemas

type UserRegisterBody struct {
	Name      string
	Birthdate string
	Email     string
	Password  string
}

// Not being used yet
type UserRegisterResponse struct{}
