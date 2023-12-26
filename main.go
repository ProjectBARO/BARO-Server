package main

import (
	"gdsc/baro/controllers"
	"gdsc/baro/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	router.GET("/health", controllers.HealthCheckController{}.HealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
