package auth_realm

import (
	"errors"
	"keeper/internal/auth"
	"keeper/internal/auth/apikey"
	"keeper/internal/auth/jwt"
	"keeper/internal/config"
	"keeper/internal/repository"
)

type AuthRealm struct {
	JwtService    jwt.IJwtService
	ApiKeyService apikey.IAPIKeyService
}

type IAuthRealm interface {
	Authenticate(credential *auth.Credential) (*auth.AuthResponse, error)
}

func NewAuthRealm(cfg *config.Config, apiKeyRepository repository.IAPIKeyRepository, userRepository repository.IUserRepository) IAuthRealm {
	return &AuthRealm{
		JwtService:    jwt.NewJwtService(cfg, userRepository),
		ApiKeyService: apikey.NewAPIKeyService(apiKeyRepository, userRepository),
	}
}

func (r *AuthRealm) Authenticate(credential *auth.Credential) (*auth.AuthResponse, error) {
	if credential.Type == auth.CredentialTypeAPIKey {
		return r.ApiKeyService.Authenticate(credential)
	}
	if credential.Type == auth.CredentialTypeJWT || credential.Type == auth.CredentialTypeRefreshJWT {
		return r.JwtService.Authenticate(credential)
	}
	return nil, errors.New("credential type is invalid")
}
