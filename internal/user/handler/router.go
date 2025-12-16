package handler

import (
	"github.com/gin-gonic/gin"

	"ticket-io/internal/user/service"
)

func RegisterRoutes(r *gin.RouterGroup, userService *service.Service) {
	handler := NewHandler(userService)

	r.GET("/users", handler.GetAll)
	r.GET("/users/:id", handler.GetByID)
	r.POST("/users", handler.Create)
}
