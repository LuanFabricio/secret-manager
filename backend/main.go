package main

/*
import "fmt"
import "net/http"
*/

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"secret-manager/backend/controllers"
)

func main() {
	godotenv.Load();

	// NOTE: Maybe move to New() and add the log middleware
	r := gin.Default();

	controllers.BindRoutes(r);

	r.Run();
}
