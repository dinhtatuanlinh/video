package restful

import (
	"github.com/dinhtatuanlinh/video/internal/delivery/restful/middleware"
	"github.com/dinhtatuanlinh/video/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router  *gin.Engine
	useCase usecase.UseCase
}

func NewServer(usecase usecase.UseCase) (*Server, error) {
	server := &Server{
		useCase: usecase,
	}

	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("currency", validCurrency)
	//}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // or restrict to specific domains
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(middleware.HttpLogger())
	//router.GET("/video/:name", server.)
	router.GET("/health_check", server.HealthCheckHandler)
	router.POST("/video/download", server.DownloadVideoHandler)
	router.POST("/video/create", server.CreateVideoHandler)
	router.POST("/video/category", server.CreateVideoCategoryHandler)
	router.Static("/downloads", "/videos")

	//router.POST("/admin", server.CreateAdminHandler)

	//router.GET("/operator/verify_email", server.VerifyEmailHandler)
	//router.POST("/operator/resend_verify/:username", server.useCaseVerifyEmail.ResendVerifyEmail)

	//authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))
	//authRoutes.POST("/operator", server.CreateOperatorHandler)
	//authRoutes.POST("/accounts", server.createAccount)

	server.Router = router
}

// Start starts the HTTP server on the specified address.
// This method is kept for backward compatibility with tests.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
