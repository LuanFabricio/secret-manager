package controllers

import (
	"github.com/gin-gonic/gin"

	"secret-manager/backend/controllers/user"
	auth_controller "secret-manager/backend/controllers/auth"
	auth_middleware "secret-manager/backend/middlewares/auth"
)

func BindRoutes(r *gin.Engine) {
	r.POST("/auth", auth_controller.AuthCredentials)

	authorized := r.Group("/")
	authorized.Use(auth_middleware.AuthMiddleware())
	{
		authorized.GET("/user/id/:id", user.GetUserById);
		authorized.GET("/user/username/:username", user.GetUserByUsername);
		authorized.POST("/user", user.CreateUser);
	}
}
