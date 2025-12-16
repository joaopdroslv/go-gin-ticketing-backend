package dto

type CreateUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Birthdate string `json:"birthdate" binding:"required"`
	StatusID  int64  `json:"status_id" binding:"required"`
}
