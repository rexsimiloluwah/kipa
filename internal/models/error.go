package models

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidObjectID   = errors.New("invalid object id")
	ErrUsersNotFound     = errors.New("users not found")
	ErrUpdatingUser      = errors.New("error updating user")
	ErrDeletingUser      = errors.New("error deleting user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrAPIKeyNotFound    = errors.New("api key not found")
	ErrUpdatingAPIKey    = errors.New("error updating api key")
	ErrRevokingAPIKey    = errors.New("error revoking api key")
	ErrRevokingAPIKeys   = errors.New("error revoking api keys")
	ErrDeletingAPIKey    = errors.New("error deleting api key")
	ErrDeletingAPIKeys   = errors.New("error deleting api keys")
	ErrAPIKeysNotFound   = errors.New("api keys not found")
	ErrIncorrectPassword = errors.New("incorrect password")
)
