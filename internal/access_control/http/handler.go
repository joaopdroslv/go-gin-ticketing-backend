package achttp

import (
	"go-gin-ticketing-backend/internal/access_control/schemas"
	"go-gin-ticketing-backend/internal/access_control/service"
	"go-gin-ticketing-backend/internal/shared/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *acservice.PermissionService
}

func NewPermissionHandler(permissionService *acservice.PermissionService) *PermissionHandler {

	return &PermissionHandler{permissionService: permissionService}
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {

	var query acschemas.GetAllPermissionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		schemas.Failed(c, http.StatusBadRequest, "invalid query params")
		return
	}
	query.NormalizePagination()

	h.permissionService.GetAllPermissions(c, &query)
	// response, err := h.permissionService.GetAllPermissions(c, &query)
	// if err != nil {
	// 	schemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
	// 	return
	// }

	schemas.OK(c, gin.H{"message": "Ok"})
	// schemas.OK(c, &response)
}
