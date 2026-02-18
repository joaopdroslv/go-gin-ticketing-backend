package schemas

type UserCreateBody struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Birthdate string `json:"birthdate" binding:"required"`
	StatusID  int64  `json:"status_id" binding:"required"`
}

type UserUpdateBody struct {
	Name      *string `json:"name" binding:"omitempty,min=2"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Birthdate *string `json:"birthdate" binding:"omitempty,email"`
}
