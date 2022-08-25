package server

import (
	"fmt"
	"keeper/internal/config"
	"keeper/internal/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Server      *echo.Echo
	Cfg         *config.Config
	Handler     *handlers.Handler
	Middlewares *Middleware
}

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

	handler := handlers.NewHandler(dbClient)
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
