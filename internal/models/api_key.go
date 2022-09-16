package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// API Key struct
type APIKey struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID    `bson:"user_id,omitempty" json:"user_id"`
	MaskID      string                `bson:"mask_id,omitempty" json:"-"`
	Name        string                `bson:"name,omitempty" json:"name"`
	Salt        string                `bson:"salt,omitempty" json:"-"`
	Hash        string                `bson:"hash,omitempty" json:"-"`
	Revoked     bool                  `bson:"revoked,omitempty" json:"revoked"`
	Key         string                `bson:"key,omitempty" json:"key,omitempty"`
	KeyType     string                `bson:"key_type,omitempty" json:"key_type"`
	Role        string                `bson:"role,omitempty" json:"role"`
	Permissions APIKeyPermissionsList `bson:"permissions,omitempty" json:"permissions"`
	ExpiresAt   primitive.DateTime    `bson:"expires_at,omitempty" json:"expires_at"`
	CreatedAt   primitive.DateTime    `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   primitive.DateTime    `bson:"updated_at,omitempty" json:"updated_at"`
}
