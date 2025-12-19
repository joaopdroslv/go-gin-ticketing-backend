package middleware

import (
	"strconv"
	"ticket-io/internal/auth/service"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(accessControl service.AccessControl, permission string) gin.HandlerFunc {

	return func(c *gin.Context) {

		userIDStr := c.GetHeader("user_id") // From JWT authentication

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid user_id"})
			return
		}

		// Skipping scope validation for now

		allowed, err := accessControl.ValidateUserPermission(c.Request.Context(), int64(userID), permission)
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
