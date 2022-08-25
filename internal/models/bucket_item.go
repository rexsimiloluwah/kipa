package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Bucket item struct
type BucketItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	BucketID  primitive.ObjectID `bson:"bucket_id,omitempty" json:"bucket_id"`
	BucketUID string             `bson:"bucket_uid,omitempty" json:"bucket_uid"`
	Key       string             `bson:"key,omitempty" json:"key"`
	Data      interface{}        `bson:"data,omitempty" json:"data"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at"`
}
