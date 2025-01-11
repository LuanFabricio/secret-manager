package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"secret-manager/backend/models/user"
	"secret-manager/backend/services/database"
)

var db database.Database = database.GetConnection()

func CreateUser(c *gin.Context) {
	var user_data user.UserDTO;

	err := c.ShouldBindJSON(&user_data);
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": err.Error(),
		});
		return;
	}

	if user_data.Username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": "Error: Username not provided",
		});
		return;
	}

	if user_data.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": "Error: Password not provided",
		});
		return;
	}

	new_user, err := user.Create(db, user_data.Username, user_data.Password);
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H {
			"message": err.Error(),
		});
		return;
	}

	c.IndentedJSON(http.StatusCreated, gin.H {
		"id": new_user.ID,
		"username": new_user.Username,
	});
}

func GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H {
			"message": err.Error(),
		});
		return;
	}

	find_user, err := user.FindByID(db, uint(id));
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H {
			"message": err.Error(),
		});
		return;
	}

	c.IndentedJSON(http.StatusOK, find_user.ToH());
}

func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	find_user, err := user.FindByUsername(db, username);
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H {
			"message": err.Error(),
		});
		return;
	}

	c.IndentedJSON(http.StatusOK, find_user.ToH());
}
