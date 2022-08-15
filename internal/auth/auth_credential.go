package auth

import (
	"errors"
	"keeper/internal/models"
)

var (
	ErrCredentialNotFound = errors.New("credential not found")
)

type Credential struct {
	Type   CredentialType `json:"type"`
	JWT    string         `json:"jwt"`
	APIKey string         `json:"api_key"`
}

type CredentialType string

type AuthResponse struct {
	AuthMode   CredentialType `json:"-"` // mode of authentication i.e. jwt or api_key
	Credential Credential     `json:"credential"`
	User       *models.User   `json:"user"`
}

const (
	CredentialTypeAPIKey     = CredentialType("api_key")
	CredentialTypeJWT        = CredentialType("jwt")
	CredentialTypeRefreshJWT = CredentialType("x-refresh-token")
)

func (c CredentialType) String() string {
	return string(c)
}
