package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitHealthCheckRoute(s *Server) {
	s.Server.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is healthy!")
	})
}

func InitAuthRoutes(s *Server) {
	authRoutes := s.Server.Group("api/v1/auth")
	protectedAuthRoutes := authRoutes.Group("")
	{
		protectedAuthRoutes.Use(s.Middlewares.RequireAuth)
		protectedAuthRoutes.GET("/user", s.Handler.AuthHandler.GetAuthUser)
		protectedAuthRoutes.POST("/refresh-token", s.Handler.AuthHandler.RefreshToken)
	}
	authRoutes.POST("/register", s.Handler.UserHandler.RegisterUser)
	authRoutes.POST("/login", s.Handler.AuthHandler.Login)
}

func InitUserRoutes(s *Server) {
	userRoutes := s.Server.Group("api/v1/users")
	userProtectedRoutes := s.Server.Group("api/v1/users", s.Middlewares.RequireAuth)
	userRoutes.GET("/:userId", s.Handler.UserHandler.GetUserByID)
	userRoutes.GET("", s.Handler.UserHandler.GetAllUsers)
	userProtectedRoutes.PUT("", s.Handler.UserHandler.UpdateUser)
	userProtectedRoutes.DELETE("", s.Handler.UserHandler.DeleteUser)
}

func InitAPIKeyRoutes(s *Server) {
	apiKeyRoutes := s.Server.Group("/api/v1/api_keys")
	{
		apiKeyRoutes.Use(s.Middlewares.RequireAuth)
		apiKeyRoutes.POST("", s.Handler.APIKeyHandler.CreateAPIKey)
		apiKeyRoutes.GET("/:apiKeyId", s.Handler.APIKeyHandler.FindAPIKeyByID)
		apiKeyRoutes.PUT("/:apiKeyId", s.Handler.APIKeyHandler.UpdateAPIKey)
		apiKeyRoutes.PUT("/revoke", s.Handler.APIKeyHandler.RevokeAPIKeys)
		apiKeyRoutes.DELETE("", s.Handler.APIKeyHandler.DeleteAPIKeys)
	}
}