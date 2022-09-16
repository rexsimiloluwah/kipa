package dto

import (
	"keeper/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAPIKeyInputDTO struct {
	Name        string                       `json:"name" form:"name" validate:"required" swaggertype:"string" example:"test-key"`
	KeyType     string                       `json:"key_type" form:"key_type" validate:"max=150" swaggertype:"string" example:""`
	Role        string                       `json:"role" form:"role" validate:"max=150" swaggertype:"string" example:""`
	Permissions models.APIKeyPermissionsList `json:"permissions" form:"permissions" swaggertype:"array,string" example:""`
	ExpiresAt   *time.Time                   `json:"expires_at" form:"expires_at" swaggertype:"string" example:"2022-09-30T15:04:05-07:00"`
}

type UpdateAPIKeyInputDTO struct {
	Name        string                       `json:"name" form:"name" validate:"required" swaggertype:"string" example:"updated-test-key"`
	KeyType     string                       `json:"key_type" form:"key_type" validate:"max=150" swaggertype:"string" example:""`
	Role        string                       `json:"role" form:"role" validate:"max=150" swaggertype:"string" example:""`
	Permissions models.APIKeyPermissionsList `json:"permissions" form:"permissions" swaggertype:"array,string" example:""`
	ExpiresAt   *time.Time                   `json:"expires_at" form:"expires_at" swaggertype:"primitive,string" example:"2022-09-30T15:04:05-07:00"`
}

type CreateAPIKeyOutputDTO struct {
	ID          primitive.ObjectID           `json:"id"`
	Name        string                       `json:"name"`
	Key         string                       `json:"key"`
	Permissions models.APIKeyPermissionsList `json:"permissions"`
	ExpiresAt   primitive.DateTime           `json:"expires_at"`
	CreatedAt   primitive.DateTime           `json:"created_at"`
}

type APIKeysIDsInputDTO struct {
	Ids []string `json:"ids" validate:"min=1" form:"ids" swaggertype:"array,string"`
}
