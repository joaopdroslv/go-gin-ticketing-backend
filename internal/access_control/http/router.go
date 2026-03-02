package accesscontrol

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccessControlRoutes(r *gin.RouterGroup, handler *PermissionHandler) {

	r.GET("/permission", handler.GetAllPermissions)
}
