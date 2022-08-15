package apikey

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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
	apiKeyRepository repository.IAPIKeyRepository
	userRepository   repository.IUserRepository
}

func NewAPIKeyService(apiKeyRepository repository.IAPIKeyRepository, userRepository repository.IUserRepository) IAPIKeyService {
	return &APIKeyService{
		apiKeyRepository: apiKeyRepository,
		userRepository:   userRepository,
	}
}

// Authenticate an API Key
func (a *APIKeyService) Authenticate(credential *auth.Credential) (*auth.AuthResponse, error) {
	if credential.Type != auth.CredentialTypeAPIKey {
		return nil, errors.New("credential must be of api key type")
	}

	key := credential.APIKey
	keySplit := strings.Split(key, utils.APIKeySeperator)
	// validate that the length = 3
	if len(keySplit) != 3 {
		return nil, errors.New("invalid api key")
	}

	// obtain the mask ID
	maskID := keySplit[1]
	// find the API Key
	apiKey, err := a.apiKeyRepository.FindAPIKeyByMaskID(maskID)
	if err != nil {
		return nil, fmt.Errorf("api key is incorrect: %s", err.Error())
	}

	// decode the API Key
	decodedKey, err := base64.URLEncoding.DecodeString(apiKey.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode api key hash: %s", err.Error())
	}

	// compute hash and compare
	dk := pbkdf2.Key([]byte(credential.APIKey), []byte(apiKey.Salt), 4096, 32, sha256.New)
	if !bytes.Equal(dk, decodedKey) {
		// mismatch
		return nil, errors.New("invalid api key")
	}

	// current time > apiKey expires at time? api key is expired
	if apiKey.ExpiresAt != 0 && time.Now().After(apiKey.ExpiresAt.Time()) {
		return nil, errors.New("api key is expired")
	}

	// check if the api key has been revoked or not
	if apiKey.Revoked {
		return nil, errors.New("api key has been revoked")
	}

	// get the user payload
	user, err := a.userRepository.FindUserById(apiKey.UserID.Hex())
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	// generate the auth response
	authResponse := &auth.AuthResponse{
		AuthMode:   auth.CredentialTypeAPIKey,
		Credential: *credential,
		User:       user,
	}

	return authResponse, nil
}
