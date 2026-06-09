package user

import "go-gin-ticketing-backend/internal/shared/schemas"

type CreateUserBody struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type UpdateUserBody struct {
	Name      *string `json:"name" binding:"omitempty,min=2"`
	Birthdate *string `json:"birthdate" binding:"required"`
	Email     *string `json:"email" binding:"omitempty,email"`
}

type ResponseUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Birthdate  string `json:"birthdate"`
	Email      string `json:"email"`
	UserStatus string `json:"user_status"`
}

type GetAllUsersResponse struct {
	Items      []ResponseUser                   `json:"items"`
	Pagination schemas.ResponsePagination `json:"pagination"`
}

type DeleteUserResponse struct {
	ID      int64 `json:"id"`
	Deleted bool  `json:"deleted"`
}
