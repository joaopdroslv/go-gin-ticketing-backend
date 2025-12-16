package dto

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	Status    string `json:"status"`
}

type GetAllResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}
