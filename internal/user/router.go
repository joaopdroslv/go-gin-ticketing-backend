package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, userService *Service) {
	handler := NewHandler(userService)

	r.GET("/users", handler.GetAll)
	r.GET("/users/:id", handler.GetByID)
	r.POST("/users", handler.Create)
}
