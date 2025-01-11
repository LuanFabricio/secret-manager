package controllers

import (
	"github.com/gin-gonic/gin"

	"secret-manager/backend/controllers/user"
)

func BindRoutes(r *gin.Engine) {
	r.GET("/user/id/:id", user.GetUserById);
	r.GET("/user/username/:username", user.GetUserByUsername);
	r.POST("/user", user.CreateUser);
}
