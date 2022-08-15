package repository

import (
	"context"
	"keeper/internal/config"
	"keeper/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	bucketItemCollectionName = "bucketitems"
)

type IBucketItemRepository interface {
	ListBucketItems(bucketID string) ([]*models.BucketItem, error)
	CreateBucketItem(bucketItem *models.BucketItem) error
	UpdateBucketItem(bucketItem *models.BucketItem) error
	FindBucketItemByID(bucketItemID string) (*models.BucketItem, error)
	DeleteBucketItem(bucketItemID string) error
}

type BucketItemRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewBucketItemRepository(cfg *config.Config, dbClient *mongo.Client) IBucketItemRepository {
	bucketItemCollection := dbClient.Database(cfg.DbName).Collection(bucketItemCollectionName)
	return &BucketItemRepository{
		collection: bucketItemCollection,
		ctx:        context.TODO(),
	}
}

func (r *BucketItemRepository) ListBucketItems(bucketID string) ([]*models.BucketItem, error) {
	panic("not implemented") // TODO: Implement
}

func (r *BucketItemRepository) CreateBucketItem(bucketItem *models.BucketItem) error {
	panic("not implemented") // TODO: Implement
}

func (r *BucketItemRepository) UpdateBucketItem(bucketItem *models.BucketItem) error {
	panic("not implemented") // TODO: Implement
}

func (r *BucketItemRepository) FindBucketItemByID(bucketItemID string) (*models.BucketItem, error) {
	panic("not implemented") // TODO: Implement
}

func (r *BucketItemRepository) DeleteBucketItem(bucketItemID string) error {
	panic("not implemented") // TODO: Implement
}
