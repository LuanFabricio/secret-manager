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

func FindSecretsByUserID(c *gin.Context) {
	token := c.GetHeader("token")

	token_user_id, err := auth.ExtractTokenId(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user_id, err := strconv.Atoi(token_user_id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	secrets, err := secret.FindByUserID(db, uint(user_id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, secrets)
}

func DeleteSecretByID(c *gin.Context) {
	token := c.GetHeader("token")

	token_user_id, err := auth.ExtractTokenId(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user_id, err := strconv.Atoi(token_user_id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	secret_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	secret_to_delete, err := secret.FindByID(db, uint(secret_id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if uint(user_id) != secret_to_delete.UserID {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Only the owner can exclude the secret",
		})
		return
	}

	deleted_secret, err := secret.DeleteById(db, uint(secret_id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, deleted_secret)
}

// TODO: Improve design
func UpdateSecretByID(c* gin.Context) {
	token := c.GetHeader("token")

	token_user_id, err := auth.ExtractTokenId(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user_id, err := strconv.Atoi(token_user_id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var secret_update secret.SecretDB
	err = c.ShouldBindJSON(&secret_update)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if secret_update.UserID != uint(user_id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "The secret user id should be the token user id",
		})
		return
	}

	secret_to_update, err := secret.FindByID(db, secret_update.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if secret_to_update.UserID != uint(user_id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "The secret user id dont match with the token user id",
		})
		return
	}

	updated_secret, err := secret.UpdateByID(db, secret_update)
	if secret_to_update.UserID != uint(user_id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, updated_secret)
}
