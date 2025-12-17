package dto

type UserUpdateBody struct {
	Name      *string `json:"name"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Birthdate *string `json:"birthdate"`
}
