package user

import (
	"go-gin-ticketing-backend/internal/access_control/service"
	"go-gin-ticketing-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(
	r *gin.RouterGroup,
	handler *UserHandler,
	accessControl acservice.AccessControl,
) {

	r.GET(
		"",
		middlewares.PermissionMiddleware(accessControl, "user:list"),
		handler.GetAllUsers,
	)
	r.GET(
		"/:id",
		middlewares.PermissionMiddleware(accessControl, "user:read"),
		handler.GetUserByID,
	)
	r.POST(
		"",
		middlewares.PermissionMiddleware(accessControl, "user:create"),
		handler.CreateUser,
	)
	r.PUT(
		"/:id",
		middlewares.PermissionMiddleware(accessControl, "user:update"),
		handler.UpdateUserByID,
	)
	r.DELETE(
		"/:id",
		middlewares.PermissionMiddleware(accessControl, "user:delete"),
		handler.DeleteUserByID,
	)
}
