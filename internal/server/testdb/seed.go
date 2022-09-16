package testdb

import (
	"crypto/sha256"
	"encoding/base64"
	"keeper/internal/config"
	"keeper/internal/models"
	"keeper/internal/repository"
	"time"

	"keeper/internal/utils"

	"github.com/ddosify/go-faker/faker"
	"github.com/sirupsen/logrus"
	"github.com/xdg-go/pbkdf2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// create a random user for integration tests purposes
func SeedUser(dbClient *mongo.Client, cfg *config.Config) (*models.User, error) {
	faker := faker.NewFaker()
	userRepo := repository.NewUserRepository(cfg, dbClient)
	hashedPassword, err := utils.HashPassword("secret")
	if err != nil {
		logrus.WithError(err).Fatal("error hashing password")
		return &models.User{}, err
	}
	// create a new user
	newUser := &models.User{
		ID:        primitive.NewObjectID(),
		Firstname: faker.RandomPersonFirstName(),
		Lastname:  faker.RandomPersonLastName(),
		Email:     faker.RandomEmail(),
		Password:  hashedPassword,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	err = userRepo.CreateUser(newUser)
	if err != nil {
		logrus.WithError(err).Fatal("error seeding new user")
		return &models.User{}, err
	}
	return newUser, nil
}

// create a random api key for testing purposes
func SeedAPIKey(dbClient *mongo.Client, cfg *config.Config, userID primitive.ObjectID) (*models.APIKey, string, error) {
	faker := faker.NewFaker()
	apiKeyRepo := repository.NewAPIKeyRepository(cfg, dbClient)

	// Generate mask and key
	maskID, key := utils.GenerateAPIKey()
	// Generate a salt secret
	salt, err := utils.GenerateSecret()
	if err != nil {
		logrus.WithError(err).Error("error generating salt secret")
		return &models.APIKey{}, "", err
	}
	// Generate a hash of the key using the salt as a secret
	dk := pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha256.New)
	encodedKey := base64.URLEncoding.EncodeToString(dk)
	apiKey := &models.APIKey{
		Name:      faker.RandomProductName(),
		KeyType:   "seed",
		Role:      "seed",
		UserID:    userID,
		MaskID:    maskID,
		Hash:      encodedKey,
		Salt:      salt,
		Key:       key,
		ExpiresAt: primitive.NewDateTimeFromTime(time.Now().Add(time.Hour)),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	_, err = apiKeyRepo.CreateAPIKey(apiKey)
	if err != nil {
		logrus.WithError(err).Error("error seeding new api key")
		return &models.APIKey{}, "", err
	}
	return apiKey, key, nil
}

// create a random bucket for testing purposes
func SeedBucket(dbClient *mongo.Client, cfg *config.Config, userID primitive.ObjectID) (*models.Bucket, error) {
	faker := faker.NewFaker()
	bucketRepo := repository.NewBucketRepository(cfg, dbClient)

	newBucket := &models.Bucket{
		ID:          primitive.NewObjectID(),
		UID:         faker.RandomUUID().String(),
		UserID:      userID,
		Name:        faker.RandomProductName(),
		Description: "a seed bucket",
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := bucketRepo.CreateBucket(newBucket)
	if err != nil {
		logrus.WithError(err).Error("error seeding new bucket")
		return &models.Bucket{}, err
	}

	return newBucket, nil
}

// creating a random bucket item for testing purposes
func SeedBucketItem(dbClient *mongo.Client, cfg *config.Config, userID primitive.ObjectID, bucketID primitive.ObjectID, bucketUID string, data interface{}) (*models.BucketItem, error) {
	faker := faker.NewFaker()
	bucketItemRepo := repository.NewBucketItemRepository(cfg, dbClient)

	newBucketItem := &models.BucketItem{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		BucketID:  bucketID,
		BucketUID: bucketUID,
		Key:       faker.RandomUUID().String(),
		Data:      data,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := bucketItemRepo.CreateBucketItem(newBucketItem)
	if err != nil {
		logrus.WithError(err).Error("failed to seed new bucket item")
		return &models.BucketItem{}, err
	}

	return newBucketItem, nil
}
