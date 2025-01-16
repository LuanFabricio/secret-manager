package auth

import (
	"net/http"
	"secret-manager/backend/services/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.GetHeader("token")

		if userToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Undefined token. The token should be on header of the request",
			})
			c.Abort()
			return
		}

		if !auth.ValidateToken(userToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
