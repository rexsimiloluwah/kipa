package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	apiKeyCollectionName = "apikeys"
)

type APIKeyRepository struct {
	ctx        context.Context
	collection *mongo.Collection
}

type IAPIKeyRepository interface {
	CreateAPIKey(apiKey *models.APIKey) error
	UpdateAPIKey(apiKey *models.APIKey) error
	FindAPIKeyByID(apiKeyID string) (*models.APIKey, error)
	FindAPIKeyByMaskID(maskID string) (*models.APIKey, error)
	FindAPIKeyByHash(hash string) (*models.APIKey, error)
	FindUserAPIKeys(userID string) ([]models.APIKey, error)
	RevokeAPIKey(apiKeyID string) error
	RevokeAPIKeys(apiKeyIDs []string) error
	DeleteAPIKey(apiKeyID string) error
	DeleteAPIKeys(apiKeyIDs []string) error
}

var apiKeyDetailsProjection = bson.D{
	primitive.E{Key: "_id", Value: 1},
	primitive.E{Key: "name", Value: 1},
	primitive.E{Key: "key_type", Value: 1},
	primitive.E{Key: "user_id", Value: 1},
	primitive.E{Key: "role", Value: 1},
	primitive.E{Key: "key", Value: 1},
	primitive.E{Key: "revoked", Value: 1},
	primitive.E{Key: "expires_at", Value: 1},
	primitive.E{Key: "created_at", Value: 1},
	primitive.E{Key: "updated_at", Value: 1},
}

func NewAPIKeyRepository(cfg *config.Config, dbClient *mongo.Client) IAPIKeyRepository {
	apiKeyCollection := dbClient.Database(cfg.DbName).Collection(apiKeyCollectionName)
	return &APIKeyRepository{
		collection: apiKeyCollection,
		ctx:        context.TODO(),
	}
}

func mapIDsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	return objectIDs, nil
}

func (r *APIKeyRepository) CreateAPIKey(apiKey *models.APIKey) error {
	_, err := r.collection.InsertOne(r.ctx, apiKey)
	if err != nil {
		logrus.WithError(err).Error("error creating api key")
		return fmt.Errorf("error creating api key: %s", err.Error())
	}
	return nil
}

// Returns a single API Key by a specific key name i.e. _id, mask_id, hash etc.
func (r *APIKeyRepository) FindAPIKeyByKeyName(key string, value string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	ID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: key, Value: ID}}
	opts := options.FindOne().SetProjection(apiKeyDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(apiKey); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrAPIKeyNotFound
		}
		return nil, err
	}
	logrus.Info("found API Key: ", apiKey)
	return apiKey, nil
}

func (r *APIKeyRepository) FindAPIKeyByID(id string) (*models.APIKey, error) {
	return r.FindAPIKeyByKeyName("_id", id)
}

func (r *APIKeyRepository) FindAPIKeyByMaskID(maskID string) (*models.APIKey, error) {
	return r.FindAPIKeyByKeyName("mask_id", maskID)
}

func (r *APIKeyRepository) FindAPIKeyByHash(hash string) (*models.APIKey, error) {
	return r.FindAPIKeyByKeyName("hash", hash)
}

// Returns a user's api keys
func (r *APIKeyRepository) FindUserAPIKeys(userID string) ([]models.APIKey, error) {
	apiKeys := []models.APIKey{}
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: ID}}
	opts := options.Find().SetProjection(apiKeyDetailsProjection)
	cursor, err := r.collection.Find(r.ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Errorf("cannot find many")
		return nil, models.ErrAPIKeysNotFound
	}
	if err = cursor.All(r.ctx, &apiKeys); err != nil {
		return nil, models.ErrAPIKeysNotFound
	}
	logrus.Debug("user api keys: ", apiKeys)
	return apiKeys, nil
}

func (r *APIKeyRepository) UpdateAPIKey(apiKey *models.APIKey) error {
	filter := bson.D{primitive.E{Key: "_id", Value: apiKey.ID}}
	apiKeyByte, err := bson.Marshal(apiKey)
	if err != nil {
		return errors.New("error marshalling api key")
	}
	// create update query
	update := make(map[string]interface{}, 0)
	if err = bson.Unmarshal(apiKeyByte, update); err != nil {
		return errors.New("error unmarshalling update data")
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(r.ctx, filter, bson.D{primitive.E{Key: "$set", Value: update}}, opts)
	if err != nil {
		logrus.WithError(err).Error("error updating api key")
		return err
	}
	if result.UpsertedCount == 0 {
		return models.ErrUpdatingUser
	}
	return nil
}

func (r *APIKeyRepository) RevokeAPIKey(apiKeyID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: apiKeyID}}
	// create update query
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "revoked", Value: true}}}}
	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(r.ctx, filter, update, opts)
	if err != nil {
		logrus.WithError(err).Error("error revoking api key")
		return err
	}
	if result.UpsertedCount == 0 {
		return models.ErrRevokingAPIKey
	}
	return nil
}

func (r *APIKeyRepository) RevokeAPIKeys(apiKeyIDs []string) error {
	objectIDs, err := mapIDsToObjectIDs(apiKeyIDs)
	if err != nil {
		return errors.New("invalid api key object id")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: primitive.E{Key: "$in", Value: objectIDs}}}
	// create update query
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "revoked", Value: true}}}}
	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(r.ctx, filter, update, opts)
	if err != nil {
		logrus.WithError(err).Error("error revoking api keys")
		return models.ErrRevokingAPIKeys
	}
	if result.UpsertedCount == 0 {
		return models.ErrRevokingAPIKeys
	}
	return nil
}

func (r *APIKeyRepository) DeleteAPIKey(apiKeyID string) error {
	ID, err := primitive.ObjectIDFromHex(apiKeyID)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting api key")
		return models.ErrDeletingAPIKey
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingAPIKey
	}
	return nil
}

func (r *APIKeyRepository) DeleteAPIKeys(apiKeyIDs []string) error {
	objectIDs, err := mapIDsToObjectIDs(apiKeyIDs)
	if err != nil {
		return errors.New("invalid api key object id")
	}
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: primitive.E{Key: "$in", Value: objectIDs}}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting api keys")
		return models.ErrDeletingAPIKeys
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingAPIKeys
	}
	return nil
}
