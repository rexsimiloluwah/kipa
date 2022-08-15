package repository

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	bucketCollectionName = "buckets"
)

type IBucketRepository interface {
	CreateBucket(bucket *models.Bucket) error
	UpdateBucket(bucket *models.Bucket) error
	DeleteBucket(bucketID string) error
	FindBucketByID(bucketID string) ([]*models.Bucket, error)
	FindBucketsByUserID(userID string) ([]*models.Bucket, error)
}

type BucketRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewBucketRepository(cfg *config.Config, dbClient *mongo.Client) IBucketRepository {
	bucketCollection := dbClient.Database(cfg.DbName).Collection(bucketCollectionName)
	return &BucketRepository{
		collection: bucketCollection,
		ctx:        context.TODO(),
	}
}

func (r *BucketRepository) CreateBucket(bucket *models.Bucket) error {
	panic("not implemented") // TODO: Implement
}

func (r *BucketRepository) UpdateBucket(bucket *models.Bucket) error {
	panic("not implemented") // TODO: Implement
}

func (r *BucketRepository) DeleteBucket(bucketID string) error {
	panic("not implemented") // TODO: Implement
}

func (r *BucketRepository) FindBucketByID(bucketID string) ([]*models.Bucket, error) {
	panic("not implemented") // TODO: Implement
}

func (r *BucketRepository) FindBucketsByUserID(userID string) ([]*models.Bucket, error) {
	panic("not implemented") // TODO: Implement
}
