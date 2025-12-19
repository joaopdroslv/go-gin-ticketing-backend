package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, authHandler *UserAuthHandler, jwtSecret string) {}
