package controllers

import (
	"github.com/gin-gonic/gin"

	"secret-manager/backend/controllers/user"
)

func BindRoutes(r *gin.Engine) {
	r.GET("/user", user.GetUserById);
	r.POST("/user", user.CreateUser);
}
