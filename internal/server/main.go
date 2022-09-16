package server

import (
	"fmt"
	"keeper/internal/config"
	"keeper/internal/handlers"
	"net/http"

	_ "keeper/internal/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Server      *echo.Echo
	Cfg         *config.Config
	Handler     *handlers.Handler
	Middlewares *Middleware
}

// @title Kipa
// @version 0.1.0
// @description API Documentation for Kipa - your secure & serverless key/value store
// @termsOfService http://swagger.io/terms/

// @contact.name   Similoluwa Okunowo
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization

// @host localhost:5050
// @BasePath /api/v1
// @schemes http
func NewServer(cfg *config.Config, dbClient *mongo.Client) *Server {
	e := echo.New()
	// middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "\n[${status}]: ${method} ${uri}",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// swagger
	e.GET("/docs/*", echoSwagger.WrapHandler)

	handler := handlers.NewHandler(cfg, dbClient)
	middlewares := NewMiddleware(cfg, dbClient)

	return &Server{
		Server:      e,
		Cfg:         cfg,
		Handler:     handler,
		Middlewares: middlewares,
	}
}

func (s *Server) RegisterRoutes() {
	InitHealthCheckRoute(s)
	InitAuthRoutes(s)
	InitUserRoutes(s)
	InitAPIKeyRoutes(s)
	InitBucketRoutes(s)
	InitBucketItemRoutes(s)
	InitPublicRoutes(s)
}

// TODO: Add graceful shutdown

// Start the server
func (s *Server) Start() {
	PORT := s.Cfg.Port

	if PORT == "" {
		PORT = "1323"
	}
	s.RegisterRoutes()
	s.Server.Logger.Fatal(s.Server.Start(fmt.Sprintf(":%s", PORT)))
}
