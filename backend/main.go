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

	r := gin.Default();

	controllers.BindRoutes(r);

	r.Run();
}
