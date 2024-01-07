package main

import (
	reportController "gdsc/baro/app/report/controllers"
	reportRepository "gdsc/baro/app/report/repositories"
	reportService "gdsc/baro/app/report/services"
	userController "gdsc/baro/app/user/controllers"
	userRepository "gdsc/baro/app/user/repositories"
	userService "gdsc/baro/app/user/services"
	videoController "gdsc/baro/app/video/controllers"
	videoRepository "gdsc/baro/app/video/repositories"
	videoService "gdsc/baro/app/video/services"
	"gdsc/baro/docs"
	"gdsc/baro/global"
	"gdsc/baro/global/auth"
	"gdsc/baro/global/config"
	"gdsc/baro/global/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	UserCtrl   *userController.UserController
	ReportCtrl *reportController.ReportController
	VideoCtrl  *videoController.VideoController
	Router     *gin.Engine
}

func (app *App) Init() {
	app.DocsInit()
	app.InitRouter()
}

func (app *App) DocsInit() {
	docs.SwaggerInfo.Title = "Baro Server API"
	docs.SwaggerInfo.Description = "This is a Baro Server API Document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
}

func (app *App) InitRouter() {
	app.Router = gin.Default()

	app.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	JWT_SECRET := os.Getenv("JWT_SECRET")
	authMiddleware := auth.NewAuthentication(JWT_SECRET)

	DB, connectionErr := config.ConnectDatabase()
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}

	userRepository := userRepository.NewUserRepository(DB)
	userUtil := utils.NewUserUtil(userRepository)

	userService := userService.NewUserService(userRepository, userUtil)
	app.UserCtrl = userController.NewUserController(userService)

	reportRepository := reportRepository.NewReportRepository(DB)
	reportService := reportService.NewReportService(reportRepository, userUtil)
	app.ReportCtrl = reportController.NewReportController(reportService)

	videoRepository := videoRepository.NewVideoRepository(DB)
	videoService := videoService.NewVideoService(videoRepository)
	app.VideoCtrl = videoController.NewVideoController(videoService)

	openAPI := app.Router.Group("/")
	{
		openAPI.GET("/health", global.HealthCheckController{}.HealthCheck)
		openAPI.POST("/login", func(c *gin.Context) { app.UserCtrl.LoginOrRegisterUser(c) })
	}

	secureAPI := app.Router.Group("/")
	secureAPI.Use(authMiddleware.StripTokenMiddleware())
	{
		secureAPI.GET("/users/me", func(c *gin.Context) { app.UserCtrl.GetUserInfo(c) })
		secureAPI.PUT("/users/me", func(c *gin.Context) { app.UserCtrl.UpdateUserInfo(c) })
		secureAPI.DELETE("/users/me", func(c *gin.Context) { app.UserCtrl.DeleteUser(c) })

		secureAPI.POST("/predict", func(c *gin.Context) { app.ReportCtrl.Predict(c) })
		secureAPI.GET("/predict", func(c *gin.Context) { app.ReportCtrl.GetPredict(c) })
		secureAPI.GET("/predict/all", func(c *gin.Context) { app.ReportCtrl.GetPredicts(c) })

		secureAPI.GET("/videos", func(c *gin.Context) { app.VideoCtrl.GetVideos(c) })
	}
}

// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	app := App{}
	app.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := app.Router.Run(":" + port)
	if err != nil {
		return
	}
}
