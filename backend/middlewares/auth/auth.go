package auth

import (
	"log"
	"net/http"
	"secret-manager/backend/services/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.GetHeader("token")

		log.Printf("User token: %v\n", userToken)
		if userToken == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "Undefined token. The token should be on header of the request",
			})
			return
		}

		log.Printf("Valid token: %v\n", auth.ValidateToken(userToken))
		if !auth.ValidateToken(userToken) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			return
		}

		c.Next()
	}
}
