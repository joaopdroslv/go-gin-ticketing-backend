package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, userService *Service) {
	handler := NewUserHandler(userService)

	r.GET("/users/:id", handler.Get)
	r.POST("/users", handler.Create)
}
