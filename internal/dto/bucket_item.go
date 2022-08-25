package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateBucketItemInputDTO struct {
	Key  string      `json:"key" form:"key"`
	Data interface{} `json:"data" form:"data"`
}

type UpdateBucketItemInputDTO struct {
	Key  string      `json:"key" form:"key"`
	Data interface{} `json:"data" form:"data"`
}

type CreateBucketItemOutputDTO struct {
	ID        primitive.ObjectID `json:"id"`
	BucketUID string             `json:"bucket_uid"`
	Key       string             `json:"key"`
	Data      interface{}        `json:"data"`
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
