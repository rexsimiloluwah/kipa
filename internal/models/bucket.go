package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Bucket struct - A bucket is sort of a container that holds key/value pairs (for a specific group)
type Bucket struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UID         string                `bson:"uid,omitempty" json:"uid"`         // a shorter id for the bucket
	UserID      primitive.ObjectID    `bson:"user_id,omitempty" json:"user_id"` // bucket owner
	Name        string                `bson:"name,omitempty" json:"name"`
	Description string                `bson:"description,omitempty" json:"description"`
	Permissions BucketPermissionsList `bson:"permissions,omitempty" json:"permissions"`
	CreatedAt   primitive.DateTime    `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   primitive.DateTime    `bson:"updated_at,omitempty" json:"updated_at"`
}
