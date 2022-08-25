package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handler struct {
	UserHandler       IUserHandler
	APIKeyHandler     IAPIKeyHandler
	AuthHandler       IAuthHandler
	BucketHandler     IBucketHandler
	BucketItemHandler IBucketItemHandler
}

func NewHandler(dbClient *mongo.Client) *Handler {
	h := &Handler{
		UserHandler:       NewUserHandler(dbClient),
		APIKeyHandler:     NewAPIKeyHandler(dbClient),
		AuthHandler:       NewAuthHandler(dbClient),
		BucketHandler:     NewBucketHandler(dbClient),
		BucketItemHandler: NewBucketItemHandler(dbClient),
	}
	return h
}
