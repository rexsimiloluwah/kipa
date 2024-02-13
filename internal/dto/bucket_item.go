package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateBucketItemInputDTO struct {
	Key  string      `json:"key" form:"key" validate:"required,min=2,max=200" swaggertype:"string" example:"test-key"`
	Data interface{} `json:"data" form:"data" swaggertype:"primitive,string" example:""`
	TTL  int         `json:"ttl" form:"ttl" validate:"min=0" swaggertype:"integer"`
}

type UpdateBucketItemInputDTO struct {
	Key  string      `json:"key" form:"key" validate:"required,min=2,max=200" swaggertype:"string" example:"test-key"`
	Data interface{} `json:"data" form:"data" swaggertype:"primitive,string" example:""`
	TTL  int         `json:"ttl" form:"ttl" swaggertype:"integer" example:""`
}

type CreateBucketItemOutputDTO struct {
	ID        primitive.ObjectID `json:"id"`
	BucketUID string             `json:"bucket_uid"`
	Key       string             `json:"key"`
	Data      interface{}        `json:"data"`
	TTL       int                `json:"ttl"`
	Type      string             `json:"type"`
	CreatedAt primitive.DateTime `json:"created_at"`
}

type BucketItemDetailDTO struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	BucketUID string             `json:"bucket_uid"`
	Key       string             `json:"key"`
	Data      interface{}        `json:"data"`
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}
