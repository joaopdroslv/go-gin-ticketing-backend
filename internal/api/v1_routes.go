package api

import (
	accesscontrolhttp "go-gin-ticketing-backend/internal/access_control/http"
	"go-gin-ticketing-backend/internal/auth"
	"go-gin-ticketing-backend/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterV1(apiGroup *gin.RouterGroup, dependencies Dependencies) {

	v1Group := apiGroup.Group("/v1")

	authGroup := v1Group.Group("/auth")
	auth.RegisterAuthRoutes(authGroup, dependencies.AuthHandler)

	userGroup := v1Group.Group("/users")
	userGroup.Use(*dependencies.JWTMiddleware)
	user.RegisterUserRoutes(userGroup, dependencies.UserHandler, dependencies.PermissionService)

	// NOTE: no auth middleware for now
	accessControlGroup := v1Group.Group("/access-control")
	accesscontrolhttp.RegisterAccessControlRoutes(accessControlGroup, dependencies.PermissionHandler)
}
