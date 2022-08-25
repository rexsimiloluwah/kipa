package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAPIKeyInputDTO struct {
	Name      string     `json:"name" form:"name"`
	KeyType   string     `json:"key_type" form:"key_type"`
	Role      string     `json:"role" form:"role"`
	ExpiresAt *time.Time `json:"expires_at" form:"expires_at"`
}

type UpdateAPIKeyInputDTO struct {
	Name      string     `json:"name" form:"name"`
	KeyType   string     `json:"key_type" form:"key_type"`
	Role      string     `json:"role" form:"role"`
	ExpiresAt *time.Time `json:"expires_at" form:"expires_at"`
}

type CreateAPIKeyOutputDTO struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Key       string             `json:"key"`
	ExpiresAt primitive.DateTime `json:"expires_at"`
	CreatedAt primitive.DateTime `json:"created_at"`
}

type APIKeysIDsInputDTO struct {
	Ids []string `json:"ids"`
}
