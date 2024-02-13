package handlers

import (
	"keeper/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	UserHandler         IUserHandler
	APIKeyHandler       IAPIKeyHandler
	AuthHandler         IAuthHandler
	BucketHandler       IBucketHandler
	BucketItemHandler   IBucketItemHandler
	PublicRoutesHandler IPublicRoutesHandler
}

func NewHandler(cfg *config.Config, dbClient *mongo.Client) *Handler {
	h := &Handler{
		UserHandler:         NewUserHandler(cfg, dbClient),
		APIKeyHandler:       NewAPIKeyHandler(cfg, dbClient),
		AuthHandler:         NewAuthHandler(cfg, dbClient),
		BucketHandler:       NewBucketHandler(cfg, dbClient),
		BucketItemHandler:   NewBucketItemHandler(cfg, dbClient),
		PublicRoutesHandler: NewPublicRoutesHandler(),
	}
	return h
}
