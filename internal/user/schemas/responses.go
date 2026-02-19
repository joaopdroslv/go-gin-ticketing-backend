package schemas

type ResponseUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Birthdate  string `json:"birthdate"`
	UserStatus string `json:"user_status"`
}

type GetAllResponse struct {
	Total int64          `json:"total"`
	Users []ResponseUser `json:"users"`
}

type UserDeleteResponse struct {
	ID      int64 `json:"id"`
	Deleted bool  `json:"deleted"`
}
