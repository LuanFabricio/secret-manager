package controllers

import (
	"github.com/gin-gonic/gin"

	auth_middleware "secret-manager/backend/middlewares/auth"
	"secret-manager/backend/controllers/user"
)

func BindRoutes(r *gin.Engine) {

	authorized := r.Group("/")
	authorized.Use(auth_middleware.AuthMiddleware())
	{
		authorized.GET("/user/id/:id", user.GetUserById);
		authorized.GET("/user/username/:username", user.GetUserByUsername);
		authorized.POST("/user", user.CreateUser);
	}
}
