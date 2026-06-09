package api

import (
	"go-gin-ticketing-backend/internal/access_control/http"
	"go-gin-ticketing-backend/internal/access_control/service"
	"go-gin-ticketing-backend/internal/auth"
	"go-gin-ticketing-backend/internal/user"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler       *auth.AuthHandler
	UserHandler       *user.UserHandler
	PermissionHandler *achttp.PermissionHandler
	JWTMiddleware     *gin.HandlerFunc
	PermissionService *acservice.PermissionService
}
