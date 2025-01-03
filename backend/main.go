package main

/*
import "fmt"
import "net/http"
*/

import (
	"github.com/gin-gonic/gin"

	"secret-manager/backend/controllers"
)

func main() {
	r := gin.Default();
	r.GET("/ping", func (ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		});
	})

	controllers.BindRoutes(r);

	r.Run();
}
