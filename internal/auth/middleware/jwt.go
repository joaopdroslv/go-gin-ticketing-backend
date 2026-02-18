package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"ticket-io/internal/auth/schemas"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(jwtSecret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&schemas.CustomClaims{},
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*schemas.CustomClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", userID)
		c.Set("role", claims.Role)
		c.Set("is_system", claims.Role == "system")

		c.Next()
	}
}
