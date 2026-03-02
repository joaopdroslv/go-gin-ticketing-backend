package accesscontrol

import (
	schemas "go-gin-ticketing-backend/internal/access_control/schemas"
	service "go-gin-ticketing-backend/internal/access_control/service"
	sharedschemas "go-gin-ticketing-backend/internal/shared/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {

	return &PermissionHandler{permissionService: permissionService}
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {

	var query schemas.GetAllPermissionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		sharedschemas.Failed(c, http.StatusBadRequest, "invalid query params")
		return
	}
	query.NormalizePagination()

	h.permissionService.GetAllPermissions(c, &query)
	// response, err := h.permissionService.GetAllPermissions(c, &query)
	// if err != nil {
	// 	sharedschemas.Failed(c, http.StatusInternalServerError, "sorry, something went wrong")
	// 	return
	// }

	sharedschemas.OK(c, gin.H{"message": "Ok"})
	// sharedschemas.OK(c, &response)
}
