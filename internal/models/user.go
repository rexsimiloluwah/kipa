package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct
type User struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Firstname            string             `bson:"firstname,omitempty" json:"firstname"`
	Lastname             string             `bson:"lastname,omitempty" json:"lastname"`
	Username             string             `bson:"username,omitempty" json:"username"`
	Email                string             `bson:"email,omitempty" json:"email"`
	Password             string             `bson:"password,omitempty" json:"-"`
	Role                 string             `bson:"role,omitempty" json:"role"`
	RegistrationProvider string             `bson:"registration_provider,omitempty" json:"registration_provider"`
	HashedRefreshToken   string             `bson:"hashed_refresh_token,omitempty" json:"hashed_refresh_token"`
	EmailVerified        bool               `bson:"email_verified,omitempty" json:"email_verified"`
	CreatedAt            primitive.DateTime `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt            primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at"`
}

func (u User) String() {
	fmt.Printf("ID: %s, Firstname: %s, Lastname: %s, Email: %s, Password: %s, RegistrationProvider: %s, HashedRefreshToken: %s, EmailVerified: %v, CreatedAt: %v, UpdatedAt: %v",
		u.ID,
		u.Firstname,
		u.Lastname,
		u.Email,
		u.Password,
		u.RegistrationProvider,
		u.HashedRefreshToken,
		u.EmailVerified,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
