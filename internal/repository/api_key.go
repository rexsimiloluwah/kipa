package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/models"
	"keeper/internal/utils"

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

var apiKeyDetailsProjection = bson.D{
	primitive.E{Key: "_id", Value: 1},
	primitive.E{Key: "name", Value: 1},
	primitive.E{Key: "key_type", Value: 1},
	primitive.E{Key: "user_id", Value: 1},
	primitive.E{Key: "role", Value: 1},
	primitive.E{Key: "key", Value: 1},
	primitive.E{Key: "mask_id", Value: 1},
	primitive.E{Key: "salt", Value: 1},
	primitive.E{Key: "hash", Value: 1},
	primitive.E{Key: "revoked", Value: 1},
	primitive.E{Key: "permissions", Value: 1},
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

// Saves a new API Key in the database
// Returns the ID of the saved API Key and an error
func (r *APIKeyRepository) CreateAPIKey(apiKey *models.APIKey) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(r.ctx, apiKey)
	if err != nil {
		logrus.WithError(err).Error("error creating api key")
		return primitive.ObjectID{}, fmt.Errorf("error creating api key: %s", err.Error())
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Utility function for finding an API Key by a specific field name
// Accepts the value and key name for the field
// Returns the found API Key or an error
func (r *APIKeyRepository) FindAPIKeyByKeyName(key string, value string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	filter := bson.D{primitive.E{Key: key, Value: value}}
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

// Find an API Key by the ID
// Accepts the API Key ID, Returns the found API Key and an error
func (r *APIKeyRepository) FindAPIKeyByID(id string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.APIKey{}, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	opts := options.FindOne().SetProjection(apiKeyDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(apiKey); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.APIKey{}, models.ErrAPIKeyNotFound
		}
		return &models.APIKey{}, err
	}
	logrus.Info("found API Key: ", apiKey)
	return apiKey, nil
}

// Find an API Key by mask ID
// Accepts the mask ID, Returns the found API Key and an error
func (r *APIKeyRepository) FindAPIKeyByMaskID(maskID string) (*models.APIKey, error) {
	return r.FindAPIKeyByKeyName("mask_id", maskID)
}

// Find an API Key by hash
// Accepts the API Key hash, Returns the found API Key and an error
func (r *APIKeyRepository) FindAPIKeyByHash(hash string) (*models.APIKey, error) {
	return r.FindAPIKeyByKeyName("hash", hash)
}

// Find a user's API Keys
// Accepts the user ID, Returns a list of the user's API Keys and an error
func (r *APIKeyRepository) FindUserAPIKeys(userID string) ([]models.APIKey, error) {
	apiKeys := []models.APIKey{}
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: ID}}
	opts := options.Find().SetProjection(apiKeyDetailsProjection).SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(r.ctx, filter, opts)

	if err != nil {
		logrus.WithError(err).Errorf("cannot find many apikeys")
		return nil, models.ErrAPIKeysNotFound
	}
	if err = cursor.All(r.ctx, &apiKeys); err != nil {
		logrus.WithError(err).Errorf("cannot find many apikeys")
		return nil, models.ErrAPIKeysNotFound
	}
	logrus.Debug("user api keys: ", apiKeys)
	return apiKeys, nil
}

func (r *APIKeyRepository) UpdateAPIKey(apiKey *models.APIKey) error {
	filter := bson.D{primitive.E{Key: "_id", Value: apiKey.ID}}
	apiKeyByte, err := bson.Marshal(apiKey)
	if err != nil {
		return errors.New("error marshalling apikey data")
	}
	// create update query
	update := make(map[string]interface{}, 0)
	if err = bson.Unmarshal(apiKeyByte, update); err != nil {
		return errors.New("error unmarshalling apikey update data")
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(r.ctx, filter, bson.D{primitive.E{Key: "$set", Value: update}}, opts)
	if err != nil {
		logrus.WithError(err).Error("error updating api key")
		return models.ErrUpdatingAPIKey
	}
	return nil
}

// Revoke an API Key
// Accepts the API Key ID and sets the 'revoked' field to true
// Returns an error on failure
func (r *APIKeyRepository) RevokeAPIKey(apiKeyID string) error {
	ID, err := primitive.ObjectIDFromHex(apiKeyID)
	if err != nil {
		return models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	// create update query
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "revoked", Value: true}}}}
	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(r.ctx, filter, update, opts)
	if err != nil {
		logrus.WithError(err).Error("error revoking api key")
		return models.ErrRevokingAPIKey
	}
	return nil
}

// Revoke multiple API Keys
// Accepts a list of API Key IDs and sets the 'revoked' status to true
// Returns an error on failure
func (r *APIKeyRepository) RevokeAPIKeys(apiKeyIDs []string) error {
	objectIDs, err := utils.MapIDsToObjectIDs(apiKeyIDs)
	if err != nil {
		return errors.New("invalid api key object id")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$in", Value: objectIDs}}}}
	// create update query
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "revoked", Value: true}}}}
	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateMany(r.ctx, filter, update, opts)
	if err != nil {
		logrus.WithError(err).Error("error revoking api keys")
		return models.ErrRevokingAPIKeys
	}
	return nil
}

// Delete a single API Key from the database
// Accepts an API Key ID, Returns an error on failure
func (r *APIKeyRepository) DeleteAPIKey(apiKeyID string) error {
	ID, err := primitive.ObjectIDFromHex(apiKeyID)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	_, err = r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting api key")
		return models.ErrDeletingAPIKey
	}
	return nil
}

// Delete multiple API Keys from the database
// Accepts a list of API Key IDs to be deleted, Returns an error on failure
func (r *APIKeyRepository) DeleteAPIKeys(apiKeyIDs []string) error {
	objectIDs, err := utils.MapIDsToObjectIDs(apiKeyIDs)
	if err != nil {
		return errors.New("invalid api key object id")
	}
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$in", Value: objectIDs}}}}
	_, err = r.collection.DeleteMany(r.ctx, filter, nil)
	if err != nil {
		logrus.WithError(err).Error("error deleting api keys")
		return models.ErrDeletingAPIKeys
	}
	return nil
}
