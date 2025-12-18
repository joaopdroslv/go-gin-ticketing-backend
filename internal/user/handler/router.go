package handler

import (
	"github.com/gin-gonic/gin"

	"ticket-io/internal/user/service/user"
)

func RegisterRoutes(r *gin.RouterGroup, userService *user.UserService) {

	userHandler := New(userService)

	r.GET("/users", userHandler.GetAll)
	r.GET("/users/:id", userHandler.GetByID)
	r.POST("/users", userHandler.Create)
	r.POST("/users/:id", userHandler.UpdateByID)
	r.DELETE("/users/:id", userHandler.DeleteByID)
}
