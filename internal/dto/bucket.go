package dto

import (
	"keeper/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBucketInputDTO struct {
	Name        string                       `json:"name" form:"name" validate:"required,min=2,max=150" swaggertype:"string" example:"test-bucket"`
	Description string                       `json:"description" form:"description" validate:"min=2,max=200" swaggertype:"string" example:""`
	Permissions models.BucketPermissionsList `json:"permissions" form:"permissions" swaggertype:"array,string" example:""`
}

type UpdateBucketInputDTO struct {
	Name        string                       `json:"name" form:"name" swaggertype:"string" example:"test-bucket"`
	Description string                       `json:"description" form:"description" swaggertype:"string" example:""`
	Permissions models.BucketPermissionsList `json:"permissions" form:"permissions" swaggertype:"array,string" example:""`
}

type CreateBucketOutputDTO struct {
	ID          primitive.ObjectID           `json:"id"`
	UID         string                       `json:"uid"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Permissions models.BucketPermissionsList `json:"permissions"`
	CreatedAt   primitive.DateTime           `json:"created_at"`
}

type BucketDetailsOutput struct {
	ID          primitive.ObjectID           `json:"id"`
	UID         string                       `json:"uid"`
	UserID      primitive.ObjectID           `json:"user_id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Permissions models.BucketPermissionsList `json:"permissions"`
	CreatedAt   primitive.DateTime           `json:"created_at"`
	UpdatedAt   primitive.DateTime           `json:"updated_at"`
	BucketItems []models.BucketItem          `json:"bucket_items"`
}
