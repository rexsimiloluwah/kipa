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
	bucketCollectionName = "buckets"
)

var bucketDetailsProjection = bson.D{
	primitive.E{Key: "_id", Value: 1},
	primitive.E{Key: "uid", Value: 1},
	primitive.E{Key: "user_id", Value: 1},
	primitive.E{Key: "name", Value: 1},
	primitive.E{Key: "description", Value: 1},
	primitive.E{Key: "permissions", Value: 1},
	primitive.E{Key: "created_at", Value: 1},
	primitive.E{Key: "updated_at", Value: 1},
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

// Save a new bucket data
func (r *BucketRepository) CreateBucket(bucket *models.Bucket) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(r.ctx, bucket)
	if err != nil {
		logrus.WithError(err).Error("error creating bucket")
		return primitive.ObjectID{}, fmt.Errorf("error creating bucket: %s", err.Error())
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Update bucket data
func (r *BucketRepository) UpdateBucket(bucket *models.Bucket) error {
	filter := bson.D{primitive.E{Key: "uid", Value: bucket.UID}} // find bucket by uid
	bucketByte, err := bson.Marshal(bucket)
	if err != nil {
		return errors.New("error marshalling bucket data")
	}
	// create update query
	update := make(map[string]interface{}, 0)
	if err = bson.Unmarshal(bucketByte, update); err != nil {
		return errors.New("error unmarshalling bucket update data")
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(r.ctx, filter, bson.D{primitive.E{Key: "$set", Value: update}}, opts)
	if err != nil {
		logrus.WithError(err).Error("error updating bucket")
		return models.ErrUpdatingBucket
	}
	return nil
}

// Returns a single bucket by id
func (r *BucketRepository) FindBucketByID(id string) (*models.Bucket, error) {
	bucket := &models.Bucket{}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	opts := options.FindOne().SetProjection(bucketDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(bucket); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrBucketNotFound
		}
		return nil, err
	}
	logrus.Info("found bucket: ", bucket)
	return bucket, nil
}

// Find a single bucket by uid
func (r *BucketRepository) FindBucketByUID(uid string) (*models.Bucket, error) {
	bucket := &models.Bucket{}
	filter := bson.D{primitive.E{Key: "uid", Value: uid}}
	opts := options.FindOne().SetProjection(bucketDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(bucket); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrBucketNotFound
		}
		return nil, err
	}
	logrus.Info("found bucket: ", bucket)
	return bucket, nil
}

// Finds an array of a user's buckets (using pagination and filtering)
func (r *BucketRepository) FindBucketsByUserIDPaged(userID string, filter bson.M, findOpts *options.FindOptions, paginationParams utils.PaginationParams) ([]models.Bucket, utils.PageInfo, error) {
	// Finds bucket items (using pagination and filtering)
	bucketItems := []models.Bucket{}

	// add the user_id to the filter object
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, utils.PageInfo{}, models.ErrInvalidObjectID
	}
	filter["user_id"] = ID

	results, pageInfo, err := utils.FindManyWithPagination(r.collection, bucketDetailsProjection, bucketItems, r.ctx, filter, findOpts, paginationParams)
	fmt.Println(results)
	if err != nil {
		return nil, utils.PageInfo{}, err
	}
	return results.([]models.Bucket), pageInfo, nil
}

// Returns an array of a user's buckets
func (r *BucketRepository) FindBucketsByUserID(userID string) ([]models.Bucket, error) {
	buckets := []models.Bucket{}
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "user_id", Value: ID}}
	opts := options.Find().SetProjection(bucketDetailsProjection).SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(r.ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Errorf("cannot find many buckets")
		return nil, models.ErrBucketsNotFound
	}
	if err = cursor.All(r.ctx, &buckets); err != nil {
		return nil, models.ErrBucketsNotFound
	}
	logrus.Debug("found buckets: ", buckets)
	return buckets, nil
}

// Delete a bucket by ID
func (r *BucketRepository) DeleteBucketByID(id string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket by id")
		return models.ErrDeletingBucket
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingBucket
	}
	return nil
}

// Delete a bucket by bucket UID
func (r *BucketRepository) DeleteBucketByUID(uid string) error {
	filter := bson.D{primitive.E{Key: "uid", Value: uid}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting bucket by uid")
		return models.ErrDeletingBucket
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingBucket
	}
	return nil
}
