package secret

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"secret-manager/backend/models/secret"
	"secret-manager/backend/services/auth"
	"secret-manager/backend/services/database"
)

var db database.Database = database.GetConnection()

func CreateSecret(c *gin.Context) {
	var secret_data secret.SecretDTO

	err := c.ShouldBindJSON(&secret_data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if secret_data.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid secret name",
		})
		return
	}

	if secret_data.Secret == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid secret",
		})
		return
	}

	secret_db, err := secret.Create(db, secret_data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusCreated, secret_db.ToH())
}

func FindSecretByID(c *gin.Context) {
	secret_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user_id_string, err := auth.ExtractTokenId(c.GetHeader("token"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user_id, err := strconv.Atoi(user_id_string)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	secret_db, err := secret.FindByID(db, uint(secret_id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if secret_db.UserID != uint(user_id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "This secret does not belongs to this user",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, secret_db.ToH())
}
