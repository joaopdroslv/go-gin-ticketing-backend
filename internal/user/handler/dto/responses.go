package dto

type ResponseUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	Status    string `json:"status"`
}

type GetAllResponse struct {
	Total int64          `json:"total"`
	Items []ResponseUser `json:"items"`
}

type UserDeleteResponse struct {
	ID      int64 `json:"id"`
	Deleted bool  `json:"deleted"`
}
