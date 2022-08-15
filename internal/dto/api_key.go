package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAPIKeyInputDTO struct {
	Name      string     `json:"name"`
	KeyType   string     `json:"key_type"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type UpdateAPIKeyInputDTO struct {
	Name      string     `json:"name"`
	KeyType   string     `json:"key_type"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type CreateAPIKeyOutputDTO struct {
	Name      string             `json:"name"`
	Key       string             `json:"key"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
