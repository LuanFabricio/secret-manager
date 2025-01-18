package secret

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"secret-manager/backend/models/secret"
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
