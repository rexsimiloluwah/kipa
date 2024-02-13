package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/models"
	"keeper/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bucketItemDetailsProjection = bson.D{
	primitive.E{Key: "_id", Value: 1},
	primitive.E{Key: "bucket_uid", Value: 1},
	primitive.E{Key: "bucket_id", Value: 1},
	primitive.E{Key: "user_id", Value: 1},
	primitive.E{Key: "key", Value: 1},
	primitive.E{Key: "data", Value: 1},
	primitive.E{Key: "ttl", Value: 1},
	primitive.E{Key: "type", Value: 1},
	primitive.E{Key: "created_at", Value: 1},
	primitive.E{Key: "updated_at", Value: 1},
}

const (
	bucketItemCollectionName = "bucketitems"
)

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

// Finds bucket items (using pagination and filtering)
func (r *BucketItemRepository) FindBucketItemsPaged(filter bson.M, opts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.BucketItem, utils.PageInfo, error) {
	bucketItems := []models.BucketItem{}
	results, pageInfo, err := utils.FindManyWithPagination(r.collection, bucketItemDetailsProjection, bucketItems, r.ctx, filter, opts, paginationParams)
	if err != nil {
		return nil, utils.PageInfo{}, err
	}
	return results.([]models.BucketItem), pageInfo, nil
}

// Finds bucket items for a specific bucket UID
// Accepts the bucket UID
// Returns the list of bucket items and an error
func (r *BucketItemRepository) FindBucketItems(bucketUID string) ([]models.BucketItem, error) {
	bucketItems := []models.BucketItem{}
	filter := bson.D{primitive.E{Key: "bucket_uid", Value: bucketUID}}
	opts := options.Find().SetProjection(bucketItemDetailsProjection).SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(r.ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Errorf("failed to find many bucket items")
		return nil, models.ErrBucketItemsNotFound
	}
	if err = cursor.All(r.ctx, &bucketItems); err != nil {
		return nil, models.ErrBucketItemsNotFound
	}
	// logrus.Debug("found bucket items: ", bucketItems)
	return bucketItems, nil
}

// Create a new bucket item
// Accepts the new bucket item data, Returns an error on failure
func (r *BucketItemRepository) CreateBucketItem(bucketItem *models.BucketItem) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(r.ctx, bucketItem)
	if err != nil {
		logrus.WithError(err).Error("error creating bucket item")
		return primitive.ObjectID{}, fmt.Errorf("error creating bucket item: %s", err.Error())
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Update a bucket item data
// Accepts the bucket item data, Returns an error on failure
func (r *BucketItemRepository) UpdateBucketItem(bucketItem *models.BucketItem, key string) error {
	filter := bson.D{
		primitive.E{Key: "bucket_uid", Value: bucketItem.BucketUID},
		primitive.E{Key: "key", Value: key},
	}
	bucketItemByte, err := bson.Marshal(bucketItem)
	if err != nil {
		return errors.New("error marshalling bucket item data")
	}
	update := make(map[string]interface{}, 0)
	if err = bson.Unmarshal(bucketItemByte, update); err != nil {
		return errors.New("error unmarshalling bucket item update data")
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(
		r.ctx,
		filter,
		bson.D{primitive.E{Key: "$set", Value: update}},
		opts,
	)
	if err != nil {
		logrus.WithError(err).Error("error updating bucket item")
		return err
	}
	return nil
}

// Increment a bucket item integer value
// Accepts the bucket UID 'bucketUID', item key 'key', and increment amount 'amount'
// Returns an error
func (r *BucketItemRepository) IncrementIntItem(bucketUID string, key string, amount int) error {
	filter := bson.D{
		primitive.E{Key: "bucket_uid", Value: bucketUID},
		primitive.E{Key: "key", Value: key},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(
		r.ctx,
		filter,
		bson.D{primitive.E{Key: "$inc", Value: bson.D{primitive.E{Key: "data", Value: amount}}}},
		opts,
	)
	if err != nil {
		logrus.WithError(err).Error("error incrementing bucket item")
		return errors.New("cannot increment non-numeric data")
	}
	return nil
}

// Find a single bucket item by the id field
// Accepts a bucket item id
// Returns the found bucket item and an error
func (r *BucketItemRepository) FindBucketItemByID(id string) (*models.BucketItem, error) {
	bucketItem := &models.BucketItem{}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	opts := options.FindOne().SetProjection(bucketItemDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(bucketItem); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrBucketItemNotFound
		}
		return nil, err
	}
	// Check if the bucket item is expired
	// the bucket item has expired if the current time is greater than the\
	// sum of the time the item was created and the specified TTL duration of the item
	if bucketItem.TTL != 0 && time.Now().After(bucketItem.CreatedAt.Time().Add(time.Duration(bucketItem.TTL))) {
		return nil, models.ErrBucketItemExpired
	}
	logrus.Info("found bucket item: ", bucketItem)
	return bucketItem, nil
}

// Find a single bucket item by the bucket UID and key name
// Accepts a bucket UID and bucket item key name
// Returns the found bucket item and an error
func (r *BucketItemRepository) FindBucketItemByKeyName(bucketUID string, key string) (*models.BucketItem, error) {
	bucketItem := &models.BucketItem{}
	filter := bson.D{primitive.E{Key: "bucket_uid", Value: bucketUID}, primitive.E{Key: "key", Value: key}}
	opts := options.FindOne().SetProjection(bucketItemDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(bucketItem); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrBucketItemNotFound
		}
		return nil, err
	}
	logrus.Info("found bucket item: ", bucketItem)
	return bucketItem, nil
}

// Delete a single bucket item based on the id field
// Accepts a bucket item id, Returns an error on failure
func (r *BucketItemRepository) DeleteBucketItemById(id string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket item by id")
		return models.ErrDeletingBucketItem
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingBucketItem
	}
	return nil
}

// Delete multiple bucket items based on the specified ids
// Accepts a list of bucket item ids, Returns an error on failure
func (r *BucketItemRepository) DeleteBucketItemsById(ids []string) error {
	objectIDs, err := utils.MapIDsToObjectIDs(ids)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$in", Value: objectIDs}}}}
	result, err := r.collection.DeleteMany(r.ctx, filter, nil)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket items")
		return models.ErrDeletingBucketItems
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingBucketItems
	}
	return nil
}

// Delete a single bucket item based on the key field
func (r *BucketItemRepository) DeleteBucketItemByKeyName(bucketUID string, key string) error {
	filter := bson.D{primitive.E{Key: "bucket_uid", Value: bucketUID}, primitive.E{Key: "key", Value: key}}
	result, err := r.collection.DeleteOne(r.ctx, filter, nil)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket item")
		return models.ErrDeletingBucketItem
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingBucketItem
	}
	return nil
}

// Deletes all the bucket items for a particular bucket
// Accepts the bucket UID, Returns an error on failure
func (r *BucketItemRepository) DeleteBucketItems(bucketUID string) error {
	filter := bson.D{primitive.E{Key: "bucket_uid", Value: bucketUID}}
	_, err := r.collection.DeleteMany(r.ctx, filter, nil)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket items")
		return models.ErrDeletingBucketItems
	}
	return nil
}
