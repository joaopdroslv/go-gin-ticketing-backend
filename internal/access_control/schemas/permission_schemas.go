package acschemas

import "go-gin-ticketing-backend/internal/shared/schemas"

type FilterPermissionQuery struct {
	Name *string `form:"name"`
}

type GetAllPermissionsQuery struct {
	FilterPermissionQuery
	schemas.PaginationQuery
}

type ResponsePermission struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetAllPermissionsResponse struct {
	Items      []ResponsePermission             `json:"items"`
	Pagination schemas.ResponsePagination `json:"pagination"`
}
