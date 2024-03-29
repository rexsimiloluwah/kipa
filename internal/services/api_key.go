package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xdg-go/pbkdf2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIKeyService struct {
	apiKeyRepo repository.IAPIKeyRepository
	Cfg        *config.Config
}

type IAPIKeyService interface {
	CreateAPIKey(data dto.CreateAPIKeyInputDTO, userID primitive.ObjectID) (dto.CreateAPIKeyOutputDTO, error)
	FindAPIKeyByID(id string) (*models.APIKey, error)
	FindUserAPIKeys(userID string) ([]models.APIKey, error)
	UpdateAPIKey(id string, data dto.UpdateAPIKeyInputDTO) error
	RevokeAPIKey(id string) error
	RevokeAPIKeys(ids []string) error
	DeleteAPIKey(id string) error
	DeleteAPIKeys(ids []string) error
}

func NewAPIKeyService(cfg *config.Config, apiKeyRepo repository.IAPIKeyRepository) IAPIKeyService {
	return &APIKeyService{
		Cfg:        cfg,
		apiKeyRepo: apiKeyRepo,
	}
}

func (a *APIKeyService) CreateAPIKey(data dto.CreateAPIKeyInputDTO, userID primitive.ObjectID) (dto.CreateAPIKeyOutputDTO, error) {
	// check expiry date of the API Key
	if data.ExpiresAt.Before(time.Now()) {
		return dto.CreateAPIKeyOutputDTO{}, errors.New("api key expires_at cannot be before now")
	}
	// Generate mask and key
	maskID, key := utils.GenerateAPIKey()
	// Generate a salt secret
	salt, err := utils.GenerateSecret()
	if err != nil {
		logrus.Errorf("error generating salt secret: %s", err.Error())
		return dto.CreateAPIKeyOutputDTO{}, err
	}
	// Generate a hash of the key using the salt as a secret
	dk := pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New)
	encodedKey := base64.URLEncoding.EncodeToString(dk)

	var permissions models.APIKeyPermissionsList
	if len(data.Permissions) == 0 {
		// set the default api key permissions list (full access)
		permissions = models.APIKEY_PERMISSIONS
	} else {
		permissions = data.Permissions
	}

	// construct new API Key
	newAPIKey := &models.APIKey{
		Name:        data.Name,
		KeyType:     data.KeyType,
		Role:        data.Role,
		UserID:      userID,
		MaskID:      maskID,
		Hash:        encodedKey,
		Salt:        salt,
		Key:         key,
		Permissions: permissions,
		ExpiresAt:   primitive.NewDateTimeFromTime(*data.ExpiresAt),
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	// save to database
	id, err := a.apiKeyRepo.CreateAPIKey(newAPIKey)
	if err != nil {
		logrus.WithError(err).Error("could not save api key to database")
		return dto.CreateAPIKeyOutputDTO{}, err
	}

	return dto.CreateAPIKeyOutputDTO{
		ID:          id,
		Name:        data.Name,
		Key:         key,
		Permissions: permissions,
		CreatedAt:   newAPIKey.CreatedAt,
		ExpiresAt:   newAPIKey.ExpiresAt,
	}, nil
}

func (a *APIKeyService) FindAPIKeyByID(id string) (*models.APIKey, error) {
	apiKey, err := a.apiKeyRepo.FindAPIKeyByID(id)
	if err != nil {
		return &models.APIKey{}, err
	}
	return apiKey, nil
}

func (a *APIKeyService) FindUserAPIKeys(userID string) ([]models.APIKey, error) {
	apiKeys, err := a.apiKeyRepo.FindUserAPIKeys(userID)
	if err != nil {
		return []models.APIKey{}, err
	}
	return apiKeys, nil
}

func (a *APIKeyService) UpdateAPIKey(id string, data dto.UpdateAPIKeyInputDTO) error {
	apiKey := &models.APIKey{
		Name:        data.Name,
		Role:        data.Role,
		KeyType:     data.KeyType,
		Permissions: data.Permissions,
	}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidObjectID
	}
	apiKey.ID = ID
	apiKey.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	if err := a.apiKeyRepo.UpdateAPIKey(apiKey); err != nil {
		return err
	}
	return nil
}

func (a *APIKeyService) RevokeAPIKey(id string) error {
	err := a.apiKeyRepo.RevokeAPIKey(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIKeyService) RevokeAPIKeys(ids []string) error {
	err := a.apiKeyRepo.RevokeAPIKeys(ids)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIKeyService) DeleteAPIKey(id string) error {
	err := a.apiKeyRepo.DeleteAPIKey(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIKeyService) DeleteAPIKeys(ids []string) error {
	err := a.apiKeyRepo.DeleteAPIKeys(ids)
	if err != nil {
		return err
	}
	return nil
}
