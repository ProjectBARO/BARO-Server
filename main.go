package main

import (
	"gdsc/baro/auth"
	"gdsc/baro/controllers"
	"gdsc/baro/models"
	"gdsc/baro/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	JWT_SECRET := os.Getenv("JWT_SECRET")
	authMiddleware := auth.NewAuthentication(JWT_SECRET)

	models.ConnectDatabase()

	userService := services.NewUserService(models.DB)
	userController := controllers.NewUserController(userService)

	openAPI := router.Group("/")
	{
		openAPI.GET("/health", controllers.HealthCheckController{}.HealthCheck)
		openAPI.POST("/login", func(c *gin.Context) {
			userController.LoginOrRegisterUser(c)
		})
	}

	secureAPI := router.Group("/")
	secureAPI.Use(authMiddleware.StripTokenMiddleware())
	{
		secureAPI.GET("/users/me", func(c *gin.Context) {
			userController.FindUserByID(c)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
