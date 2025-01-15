package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	user_model "secret-manager/backend/models/user"
	auth_service "secret-manager/backend/services/auth"
	"secret-manager/backend/services/database"
)

var db database.Database = database.GetConnection()

func AuthCredentials(c *gin.Context) {
	var user_header user_model.UserDTO

	err := c.ShouldBindJSON(&user_header)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": err.Error(),
		})
		return
	}

	id, valid := user_model.ValidateCredentials(db, user_header)
	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": "Wrong credentials",
		})
		return
	}

	token, err := auth_service.GenerateToken(*id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H {
		"token": token,
	})
}
