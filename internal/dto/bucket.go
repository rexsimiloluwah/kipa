package dto

import (
	"keeper/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBucketInputDTO struct {
	Name        string   `json:"name" form:"name"`
	Description string   `json:"description" form:"description"`
	Permissions []string `json:"permissions" form:"permissions"`
}

type UpdateBucketInputDTO struct {
	Name        string   `json:"name" form:"name"`
	Description string   `json:"description" form:"description"`
	Permissions []string `json:"permissions" form:"permissions"`
}

type CreateBucketOutputDTO struct {
	ID          primitive.ObjectID `json:"id"`
	UID         string             `json:"uid"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Permissions []string           `json:"permissions"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}

type BucketDetailsOutput struct {
	ID          primitive.ObjectID  `json:"id"`
	UID         string              `json:"uid"`
	UserID      primitive.ObjectID  `json:"user_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Permissions []string            `json:"permissions"`
	CreatedAt   primitive.DateTime  `json:"created_at"`
	UpdatedAt   primitive.DateTime  `json:"updated_at"`
	BucketItems []models.BucketItem `json:"bucket_items"`
}
