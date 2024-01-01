package main

import (
	"gdsc/baro/auth"
	"gdsc/baro/controllers"
	"gdsc/baro/docs"
	"gdsc/baro/models"
	"gdsc/baro/services"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func DocsInit() {
	docs.SwaggerInfo.Title = "Baro Server API"
	docs.SwaggerInfo.Description = "This is a Baro Server API Document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
}

// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	DocsInit()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	JWT_SECRET := os.Getenv("JWT_SECRET")
	authMiddleware := auth.NewAuthentication(JWT_SECRET)

	models.ConnectDatabase()

	userService := services.NewUserService(models.DB)
	userController := controllers.NewUserController(userService)

	openAPI := router.Group("/")
	{
		openAPI.GET("/health", controllers.HealthCheckController{}.HealthCheck)
		openAPI.POST("/login", func(c *gin.Context) { userController.LoginOrRegisterUser(c) })
	}

	secureAPI := router.Group("/")
	secureAPI.Use(authMiddleware.StripTokenMiddleware())
	{
		secureAPI.GET("/users/me", func(c *gin.Context) { userController.GetUserInfo(c) })
		secureAPI.PUT("/users/me", func(c *gin.Context) { userController.UpdateUserInfo(c) })
		secureAPI.DELETE("/users/me", func(c *gin.Context) { userController.DeleteUser(c) })
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
