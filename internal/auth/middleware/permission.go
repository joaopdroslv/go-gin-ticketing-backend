package middleware

import (
	"go-gin-ticketing-backend/internal/auth/service"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(accessControl service.AccessControl, permission string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// Bypass to any system user skip permission validation
		isSystem, _ := c.Get("is_system")
		if isSystem == true {
			c.Next()
			return
		}

		userIDAny, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		userID := userIDAny.(int64)

		// Skipping scope validation for now

		allowed, err := accessControl.HasThisPermission(c.Request.Context(), int64(userID), permission)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "authorization error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
