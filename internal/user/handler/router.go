package handler

import (
	"go-gin-ticketing-backend/internal/auth/middleware"
	"go-gin-ticketing-backend/internal/auth/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, handler *UserHandler, accessControl service.AccessControl) {

	r.GET(
		"",
		middleware.PermissionMiddleware(accessControl, "user:list"),
		handler.ListUsers,
	)
	r.GET(
		"/:id",
		middleware.PermissionMiddleware(accessControl, "user:read"),
		handler.GetUserByID,
	)
	r.POST(
		"",
		middleware.PermissionMiddleware(accessControl, "user:create"),
		handler.CreateUser,
	)
	r.PUT(
		"/:id",
		middleware.PermissionMiddleware(accessControl, "user:update"),
		handler.UpdateUserByID,
	)
	r.DELETE(
		"/:id",
		middleware.PermissionMiddleware(accessControl, "user:delete"),
		handler.DeleteUserByID,
	)
}
