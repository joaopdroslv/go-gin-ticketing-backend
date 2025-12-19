package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, userHandler *UserHandler) {

	r.GET("", userHandler.ListUsers)
	r.GET("/:id", userHandler.GetUserByID)
	r.POST("", userHandler.CreateUser)
	r.PUT("/:id", userHandler.UpdateUserByID)
	r.DELETE("/:id", userHandler.DeleteUserByID)
}
