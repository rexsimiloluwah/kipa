package apikey

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"keeper/internal/auth"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"strings"
	"time"

	"github.com/xdg-go/pbkdf2"
)

type IAPIKeyService interface {
	Authenticate(credential *auth.Credential) (*auth.AuthResponse, error)
}

type APIKeyService struct {
	apiKeyRepo repository.IAPIKeyRepository
	userRepo   repository.IUserRepository
}

func NewAPIKeyService(apiKeyRepo repository.IAPIKeyRepository, userRepo repository.IUserRepository) IAPIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

// error constants
var (
	ErrInvalidAPIKeyLength      = errors.New("invalid api key length")
	ErrAPIKeyDoesNotExist       = errors.New("api key does not exist")
	ErrFailedToDecodeAPIKeyHash = errors.New("failed to decode api key hash")
	ErrInvalidAPIKey            = errors.New("invalid api key")
	ErrExpiredAPIKey            = errors.New("api key is expired")
	ErrRevokedAPIKey            = errors.New("api key is revoked")
	ErrAPIKeyCredentialType     = errors.New("credential must be of api key type")
)

// Authenticate an API Key
func (a *APIKeyService) Authenticate(credential *auth.Credential) (*auth.AuthResponse, error) {
	if credential.Type != auth.CredentialTypeAPIKey {
		return nil, ErrAPIKeyCredentialType
	}

	key := credential.APIKey
	keySplit := strings.Split(key, utils.APIKeySeperator)
	// validate that the length = 3
	if len(keySplit) != 3 {
		return nil, ErrInvalidAPIKeyLength
	}

	// obtain the mask ID
	maskID := keySplit[1]
	// find the API Key
	apiKey, err := a.apiKeyRepo.FindAPIKeyByMaskID(maskID)
	if err != nil {
		return nil, ErrAPIKeyDoesNotExist
	}

	// decode the API Key
	decodedKey, err := base64.URLEncoding.DecodeString(apiKey.Hash)
	if err != nil {
		return nil, ErrFailedToDecodeAPIKeyHash
	}

	// compute hash and compare
	dk := pbkdf2.Key([]byte(credential.APIKey), []byte(apiKey.Salt), 4096, 32, sha256.New)
	if !bytes.Equal(dk, decodedKey) {
		// mismatch
		return nil, ErrInvalidAPIKey
	}

	// current time > apiKey expires at time? api key is expired
	if apiKey.ExpiresAt != 0 && time.Now().After(apiKey.ExpiresAt.Time()) {
		return nil, ErrExpiredAPIKey
	}

	// check if the api key has been revoked or not
	if apiKey.Revoked {
		return nil, ErrRevokedAPIKey
	}

	// get the user payload
	user, err := a.userRepo.FindUserById(apiKey.UserID.Hex())
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	// generate the auth response
	authResponse := &auth.AuthResponse{
		AuthMode:    auth.CredentialTypeAPIKey,
		Credential:  *credential,
		User:        user,
		Permissions: apiKey.Permissions,
	}

	return authResponse, nil
}
