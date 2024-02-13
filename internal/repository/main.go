package repository

import (
	"keeper/internal/models"
	"keeper/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(id string) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userId string) error
}

type IAPIKeyRepository interface {
	CreateAPIKey(apiKey *models.APIKey) (primitive.ObjectID, error)
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

type IBucketRepository interface {
	CreateBucket(bucket *models.Bucket) (primitive.ObjectID, error)
	UpdateBucket(bucket *models.Bucket) error
	DeleteBucketByID(id string) error
	DeleteBucketByUID(uid string) error
	FindBucketByID(id string) (*models.Bucket, error)
	FindBucketByUID(uid string) (*models.Bucket, error)
	FindBucketsByUserID(userID string) ([]models.Bucket, error)
	FindBucketsByUserIDPaged(userID string, filter bson.M, findOpts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.Bucket, utils.PageInfo, error)
}

type IBucketItemRepository interface {
	FindBucketItemsPaged(filter bson.M, opts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.BucketItem, utils.PageInfo, error)
	FindBucketItems(bucketUID string) ([]models.BucketItem, error)
	CreateBucketItem(bucketItem *models.BucketItem) (primitive.ObjectID, error)
	UpdateBucketItem(bucketItem *models.BucketItem, key string) error
	IncrementIntItem(bucketUID string, key string, amount int) error
	FindBucketItemByID(id string) (*models.BucketItem, error)
	FindBucketItemByKeyName(bucketUID string, key string) (*models.BucketItem, error)
	DeleteBucketItemByKeyName(bucketUID string, key string) error
	DeleteBucketItemById(id string) error
	DeleteBucketItemsById(ids []string) error
	DeleteBucketItems(bucketUID string) error
}
