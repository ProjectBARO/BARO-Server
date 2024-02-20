package main

import (
	reportController "gdsc/baro/app/report/controllers"
	reportRepository "gdsc/baro/app/report/repositories"
	reportService "gdsc/baro/app/report/services"
	userController "gdsc/baro/app/user/controllers"
	userapp "gdsc/baro/app/user/pb"
	userRepository "gdsc/baro/app/user/repositories"
	userService "gdsc/baro/app/user/services"
	videoController "gdsc/baro/app/video/controllers"
	videoapp "gdsc/baro/app/video/pb"
	videoRepository "gdsc/baro/app/video/repositories"
	videoService "gdsc/baro/app/video/services"

	userpb "gdsc/baro/protos/user"
	videopb "gdsc/baro/protos/video"

	"gdsc/baro/docs"
	"gdsc/baro/global"
	"gdsc/baro/global/auth"
	"gdsc/baro/global/config"
	"gdsc/baro/global/utils"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

type App struct {
	UserCtrl   *userController.UserController
	ReportCtrl *reportController.ReportController
	VideoCtrl  *videoController.VideoController
	Router     *gin.Engine
}

func (app *App) Init() {
	app.InitDocs()
	app.InitRouter()
}

func (app *App) InitDocs() {
	docs.SwaggerInfo.Title = "Baro Server API"
	docs.SwaggerInfo.Description = "This is a Baro Server API Document"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
}

func (app *App) RunGrpcServer(l net.Listener, userRepository userRepository.UserRepositoryInterface, userUtil utils.UserUtilInterface, videoRepository videoRepository.VideoRepositoryInterface) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryAuthInterceptor))

	userPbApp := userapp.NewUserPbApp(userRepository, userUtil)
	videoPbApp := videoapp.NewVideoPbApp(videoRepository)

	userpb.RegisterUserServiceServer(grpcServer, userPbApp)
	videopb.RegisterVideoServiceServer(grpcServer, videoPbApp)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (app *App) RunHttpServer(l net.Listener) {
	err := app.Router.RunListener(l)
	if err != nil {
		log.Fatal(err)
	}
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
		openAPI.GET("/videos", func(c *gin.Context) { app.VideoCtrl.GetVideos(c) })
		openAPI.GET("/videos/category", func(c *gin.Context) { app.VideoCtrl.GetVideosByCategory(c) })
	}

	secureAPI := app.Router.Group("/")
	secureAPI.Use(authMiddleware.StripTokenMiddleware())
	{
		secureAPI.GET("/users/me", func(c *gin.Context) { app.UserCtrl.GetUserInfo(c) })
		secureAPI.PUT("/users/me", func(c *gin.Context) { app.UserCtrl.UpdateUserInfo(c) })
		secureAPI.DELETE("/users/me", func(c *gin.Context) { app.UserCtrl.DeleteUser(c) })
		secureAPI.PUT("/users/fcm-token", func(c *gin.Context) { app.UserCtrl.UpdateFcmToken(c) })

		secureAPI.POST("/analysis", func(c *gin.Context) { app.ReportCtrl.Analysis(c) })
		secureAPI.GET("/analysis", func(c *gin.Context) { app.ReportCtrl.GetAnalysis(c) })
		secureAPI.GET("/analysis/:id", func(c *gin.Context) { app.ReportCtrl.GetAnalysisById(c) })
		secureAPI.GET("/analysis/summary", func(c *gin.Context) { app.ReportCtrl.GetAnalysisSummary(c) })
		secureAPI.GET("/analysis/all", func(c *gin.Context) { app.ReportCtrl.GetAnalyzes(c) })
		secureAPI.GET("/analysis/rank", func(c *gin.Context) { app.ReportCtrl.GetAnalysisRankAtAgeAndGender(c) })
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

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(l)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	DB, connectionErr := config.ConnectDatabase()
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}

	userRepository := userRepository.NewUserRepository(DB)
	userUtil := utils.NewUserUtil(userRepository)

	videoRepository := videoRepository.NewVideoRepository(DB)

	go app.RunGrpcServer(grpcL, userRepository, userUtil, videoRepository)
	go app.RunHttpServer(httpL)

	err = m.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
