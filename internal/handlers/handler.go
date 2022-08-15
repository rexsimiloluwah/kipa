package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handler struct {
	UserHandler   IUserHandler
	APIKeyHandler IAPIKeyHandler
	AuthHandler   IAuthHandler
}

func NewHandler(dbClient *mongo.Client) *Handler {
	h := &Handler{
		UserHandler:   NewUserHandler(dbClient),
		APIKeyHandler: NewAPIKeyHandler(dbClient),
		AuthHandler:   NewAuthHandler(dbClient),
	}
	return h
}
