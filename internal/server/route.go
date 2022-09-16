package server

import (
	"keeper/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health check routes
func InitHealthCheckRoute(s *Server) {
	s.Server.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is healthy!")
	})
}

// User authentication routes (login,register,fetch authenticated user,refresh token, forgot password etc.)
func InitAuthRoutes(s *Server) {
	authRoutes := s.Server.Group("api/v1/auth")
	protectedAuthRoutes := authRoutes.Group("")
	{
		protectedAuthRoutes.Use(s.Middlewares.RequireAuth)
		protectedAuthRoutes.GET("/user", s.Handler.AuthHandler.GetAuthUser)
	}
	authRoutes.POST("/refresh-token", s.Handler.AuthHandler.RefreshToken, s.Middlewares.RequireRefreshToken)
	authRoutes.POST("/register", s.Handler.AuthHandler.Register)
	authRoutes.POST("/login", s.Handler.AuthHandler.Login)
}

// User account management routes
func InitUserRoutes(s *Server) {
	userRoutes := s.Server.Group("api/v1/user")
	usersRoutes := s.Server.Group("api/v1/users")
	protectedUserRoutes := userRoutes.Group("")
	{
		protectedUserRoutes.Use(s.Middlewares.RequireAuth)
		protectedUserRoutes.PUT("",
			s.Handler.UserHandler.UpdateUser,
			s.Middlewares.RequireAPIKeyWriteUserPermission)
		protectedUserRoutes.PUT("/password",
			s.Handler.UserHandler.UpdateUserPassword,
			s.Middlewares.RequireAPIKeyWriteUserPermission)
		protectedUserRoutes.DELETE("",
			s.Handler.UserHandler.DeleteUser,
			s.Middlewares.RequireAPIKeyDeleteUserPermission)
	}
	usersRoutes.GET("/:userId", s.Handler.UserHandler.GetUserByID)
	usersRoutes.GET("", s.Handler.UserHandler.GetAllUsers)
}

// API Key management routes
func InitAPIKeyRoutes(s *Server) {
	apiKeysRoutes := s.Server.Group("/api/v1/api_keys")
	apiKeyRoutes := s.Server.Group("/api/v1/api_key")
	protectedAPIKeyRoutes := apiKeyRoutes.Group("")
	{
		protectedAPIKeyRoutes.Use(s.Middlewares.RequireAuth)
		protectedAPIKeyRoutes.GET("/:apiKeyId", s.Handler.APIKeyHandler.FindAPIKeyByID)
		protectedAPIKeyRoutes.POST("", s.Handler.APIKeyHandler.CreateAPIKey)
		protectedAPIKeyRoutes.PUT("/:apiKeyId", s.Handler.APIKeyHandler.UpdateAPIKey)
		protectedAPIKeyRoutes.PUT("/:apiKeyId/revoke", s.Handler.APIKeyHandler.RevokeAPIKey)
		protectedAPIKeyRoutes.DELETE("/:apiKeyId", s.Handler.APIKeyHandler.DeleteAPIKey)
	}
	protectedAPIKeysRoutes := apiKeysRoutes.Group("")
	{
		protectedAPIKeysRoutes.Use(s.Middlewares.RequireAuth)
		protectedAPIKeysRoutes.GET("", s.Handler.APIKeyHandler.FindUserAPIKeys)
		protectedAPIKeysRoutes.PUT("/revoke", s.Handler.APIKeyHandler.RevokeAPIKeys)
		protectedAPIKeysRoutes.DELETE("", s.Handler.APIKeyHandler.DeleteAPIKeys)
	}
}

// Bucket management routes
func InitBucketRoutes(s *Server) {
	bucketRoutes := s.Server.Group("/api/v1/bucket")
	bucketsRoutes := s.Server.Group("/api/v1/buckets")
	protectedBucketRoutes := bucketRoutes.Group("")
	{
		protectedBucketRoutes.Use(s.Middlewares.RequireAuth)
		protectedBucketRoutes.POST("",
			s.Handler.BucketHandler.CreateBucket,
			s.Middlewares.RequireAPIKeyBucketWritePermission,
			s.Middlewares.RequireAPIKeyBucketWritePermission)
		protectedBucketRoutes.GET("/:bucketUID",
			s.Handler.BucketHandler.FindBucketByUID,
			s.Middlewares.RequireBucketReadAccess,
			s.Middlewares.RequireAPIKeyBucketReadPermission,
		)
		protectedBucketRoutes.PUT("/:bucketUID",
			s.Handler.BucketHandler.UpdateBucket,
			s.Middlewares.RequireBucketWriteAccess,
			s.Middlewares.RequireAPIKeyBucketWritePermission,
		)
		protectedBucketRoutes.DELETE("/:bucketUID",
			s.Handler.BucketHandler.DeleteBucket,
			s.Middlewares.RequireBucketDeleteAccess,
			s.Middlewares.RequireAPIKeyBucketDeletePermission,
		)
	}
	protectedBucketsRoutes := bucketsRoutes.Group("")
	{
		protectedBucketsRoutes.Use(s.Middlewares.RequireAuth)
		protectedBucketsRoutes.GET("",
			s.Handler.BucketHandler.ListUserBuckets,
			s.Middlewares.RequireAPIKeyBucketReadPermission,
		)
	}
}

// Bucket items management routes
func InitBucketItemRoutes(s *Server) {
	bucketItemRoutes := s.Server.Group("/api/v1/item")
	bucketItemsRoutes := s.Server.Group("/api/v1/items")
	protectedBucketItemRoutes := bucketItemRoutes.Group("")
	{
		protectedBucketItemRoutes.Use(s.Middlewares.RequireAuth)
		protectedBucketItemRoutes.POST("/:bucketUID",
			s.Handler.BucketItemHandler.CreateBucketItem,
			s.Middlewares.RequireBucketItemWriteAccess,
			s.Middlewares.RequireAPIKeyWriteItemPermission,
		)
		protectedBucketItemRoutes.GET("/:bucketUID/:key",
			s.Handler.BucketItemHandler.FindBucketItemByKeyName,
			s.Middlewares.RequireBucketItemReadAccess,
			s.Middlewares.RequireAPIKeyReadItemPermission,
		)
		protectedBucketItemRoutes.PUT("/:bucketUID/:key",
			s.Handler.BucketItemHandler.UpdateBucketItemByKeyName,
			s.Middlewares.RequireBucketItemWriteAccess,
			s.Middlewares.RequireAPIKeyWriteItemPermission,
		)
		protectedBucketItemRoutes.DELETE("/:bucketUID/:key",
			s.Handler.BucketItemHandler.DeleteBucketItemByKeyName,
			s.Middlewares.RequireBucketItemDeleteAccess,
			s.Middlewares.RequireAPIKeyDeleteItemPermission,
		)
	}
	protectedBucketItemsRoutes := bucketItemsRoutes.Group("")
	{
		protectedBucketItemsRoutes.Use(s.Middlewares.RequireAuth)
		protectedBucketItemsRoutes.GET("/:bucketUID",
			s.Handler.BucketItemHandler.ListBucketItems,
			s.Middlewares.RequireBucketItemReadAccess,
			s.Middlewares.RequireAPIKeyReadItemPermission,
		)
	}
}

// Public/Miscellaneous routes
func InitPublicRoutes(s *Server) {
	publicRoutes := s.Server.Group("/api/v1/public")
	// returns an array of api key permissions
	publicRoutes.GET("/apikey-permissions", func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.APIKeyPermissions)
	})
	// returns an array of bucket-level permissions
	publicRoutes.GET("/bucket-permissions", func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.BucketPermissions)
	})
}
