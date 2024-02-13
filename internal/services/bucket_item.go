package services

import (
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BucketItemService struct {
	bucketItemRepo repository.IBucketItemRepository
	bucketRepo     repository.IBucketRepository
	cfg            *config.Config
}

type IBucketItemService interface {
	CreateBucketItem(data dto.CreateBucketItemInputDTO, userID primitive.ObjectID, bucketUID string) (*dto.CreateBucketItemOutputDTO, error)
	UpdateBucketItemByKeyName(data dto.UpdateBucketItemInputDTO, bucketUID string, key string) error
	DeleteBucketItemById(id string) error
	DeleteBucketItemsById(id []string) error
	DeleteBucketItems(bucketUID string) error
	DeleteBucketItemByKeyName(bucketUID string, key string) error
	DeleteBucketItemsByKeyName(bucketUID string, keys []string) error
	IncrementIntValue(bucketUID string, key string, amount int) error
	FindBucketItemByID(id string) (*models.BucketItem, error)
	FindBucketItemByKeyName(bucketUID string, key string) (*models.BucketItem, error)
	ListBucketItems(bucketUID string) ([]models.BucketItem, error)
	ListBucketItemsPaged(queryParams url.Values) ([]models.BucketItem, utils.PageInfo, error)
}

func NewBucketItemService(cfg *config.Config, bucketItemRepo repository.IBucketItemRepository, bucketRepo repository.IBucketRepository) IBucketItemService {
	return &BucketItemService{
		bucketItemRepo: bucketItemRepo,
		bucketRepo:     bucketRepo,
		cfg:            cfg,
	}
}

var (
	ErrKeyIsEmpty          = errors.New("bucket item key cannot be empty")
	ErrBucketItemIDIsEmpty = errors.New("bucket item id cannot be empty")
)

// Creates a new bucket item
// Accepts the bucket item input data, user ID, and bucket UID
// Returns a success response and error
func (b *BucketItemService) CreateBucketItem(data dto.CreateBucketItemInputDTO, userID primitive.ObjectID, bucketUID string) (*dto.CreateBucketItemOutputDTO, error) {
	// validation
	if utils.IsStringEmpty(data.Key) {
		return &dto.CreateBucketItemOutputDTO{}, ErrKeyIsEmpty
	}
	if utils.IsStringEmpty(bucketUID) {
		return &dto.CreateBucketItemOutputDTO{}, ErrBucketUIDIsEmpty
	}
	// check if the bucket UID exists
	bucket, err := b.bucketRepo.FindBucketByUID(bucketUID)

	if err != nil {
		return &dto.CreateBucketItemOutputDTO{}, err
	}

	// enforce unique key
	_, err = b.bucketItemRepo.FindBucketItemByKeyName(bucketUID, data.Key)

	if !errors.Is(err, models.ErrBucketItemNotFound) {
		return &dto.CreateBucketItemOutputDTO{}, fmt.Errorf("key '%s' already exists", data.Key)
	}

	fmt.Println(data.Data, utils.TypeOf(data.Data))

	// compute the data type
	dataType := utils.TypeOf(data.Data)

	newBucketItem := &models.BucketItem{
		UserID:    userID,
		BucketUID: bucketUID,
		BucketID:  bucket.ID,
		Key:       data.Key,
		Data:      data.Data,
		TTL:       data.TTL,
		Type:      dataType,
	}
	newBucketItem.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newBucketItem.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	id, err := b.bucketItemRepo.CreateBucketItem(newBucketItem)
	if err != nil {
		return &dto.CreateBucketItemOutputDTO{}, err
	}
	return &dto.CreateBucketItemOutputDTO{
		ID:        id,
		BucketUID: bucketUID,
		Key:       data.Key,
		Data:      data.Data,
		TTL:       data.TTL,
		Type:      dataType,
		CreatedAt: newBucketItem.CreatedAt,
	}, nil
}

// Update a bucket by the key name
// Accepts the update data, bucket UID and key
// Returns an error
func (b *BucketItemService) UpdateBucketItemByKeyName(data dto.UpdateBucketItemInputDTO, bucketUID string, key string) error {
	if utils.IsStringEmpty(bucketUID) {
		return ErrBucketUIDIsEmpty
	}
	if utils.IsStringEmpty(key) {
		return ErrKeyIsEmpty
	}
	updatedBucketItem := &models.BucketItem{
		BucketUID: bucketUID,
		Key:       data.Key,
		Data:      data.Data,
		TTL:       data.TTL,
	}
	err := b.bucketItemRepo.UpdateBucketItem(updatedBucketItem, key)
	if err != nil {
		return err
	}
	return nil
}

// Increments/decrements a key value by an amount
// Accepts the bucket UID, item key, and amount
// Returns an error
func (b *BucketItemService) IncrementIntValue(bucketUID string, key string, amount int) error {
	if utils.IsStringEmpty(bucketUID) {
		return ErrBucketUIDIsEmpty
	}
	if utils.IsStringEmpty(key) {
		return ErrKeyIsEmpty
	}
	err := b.bucketItemRepo.IncrementIntItem(bucketUID, key, amount)
	if err != nil {
		return err
	}
	return nil
}

// Find a bucket item by the bucket item ID
// Accepts the string value of the bucket item object ID
// Returns the found bucket item and an error
func (b *BucketItemService) FindBucketItemByID(id string) (*models.BucketItem, error) {
	if utils.IsStringEmpty(id) {
		return &models.BucketItem{}, ErrBucketItemIDIsEmpty
	}
	bucketItem, err := b.bucketItemRepo.FindBucketItemByID(id)
	if err != nil {
		return &models.BucketItem{}, err
	}
	return bucketItem, nil
}

// Find a bucket item by the key name
// Accepts the bucket UID and key name values
// Returns the found bucket item and an error
func (b *BucketItemService) FindBucketItemByKeyName(bucketUID string, key string) (*models.BucketItem, error) {
	if utils.IsStringEmpty(bucketUID) {
		return &models.BucketItem{}, ErrBucketUIDIsEmpty
	}
	if utils.IsStringEmpty(key) {
		return &models.BucketItem{}, ErrKeyIsEmpty
	}
	bucketItem, err := b.bucketItemRepo.FindBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return &models.BucketItem{}, err
	}
	return bucketItem, nil
}

// Lists all the bucket items based on the parsed query values
// Accepts the query params passed in the request object
// Returns the found bucket items and an error
func (b *BucketItemService) ListBucketItemsPaged(queryParams url.Values) ([]models.BucketItem, utils.PageInfo, error) {
	filter, findOpts, paginationParams, err := utils.ParseRequestQueryParams(queryParams)
	if err != nil {
		return []models.BucketItem{}, utils.PageInfo{}, err
	}
	bucketItems, pageInfo, err := b.bucketItemRepo.FindBucketItemsPaged(filter, findOpts, paginationParams)
	if err != nil {
		return []models.BucketItem{}, utils.PageInfo{}, err
	}
	return bucketItems, pageInfo, nil
}

// List all the bucket items for a specific bucket
// Accepts the bucket UID
// Returns the found bucket items and an error
func (b *BucketItemService) ListBucketItems(bucketUID string) ([]models.BucketItem, error) {
	if utils.IsStringEmpty(bucketUID) {
		return []models.BucketItem{}, ErrBucketUIDIsEmpty
	}
	bucketItems, err := b.bucketItemRepo.FindBucketItems(bucketUID)
	if err != nil {
		return []models.BucketItem{}, err
	}
	return bucketItems, nil
}

// Delete a single bucket item by ID
// Accepts the string value of the bucket item's object ID
// Returns an error
func (b *BucketItemService) DeleteBucketItemById(id string) error {
	if utils.IsStringEmpty(id) {
		return ErrBucketItemIDIsEmpty
	}
	err := b.bucketItemRepo.DeleteBucketItemById(id)
	if err != nil {
		return err
	}
	return nil
}

// Delete multiple bucket items
// Accepts a list of object IDs of the desired buckets
// Returns an error
func (b *BucketItemService) DeleteBucketItemsById(ids []string) error {
	err := b.bucketItemRepo.DeleteBucketItemsById(ids)
	if err != nil {
		return err
	}
	return nil
}

// Delete bucket item by key name
// Accepts the bucket UID and key name for the desired bucket item
// Returns an error
func (b *BucketItemService) DeleteBucketItemByKeyName(bucketUID string, key string) error {
	if utils.IsStringEmpty(bucketUID) {
		return ErrBucketUIDIsEmpty
	}
	if utils.IsStringEmpty(key) {
		return ErrKeyIsEmpty
	}
	err := b.bucketItemRepo.DeleteBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return err
	}
	return nil
}

// Delete bucket items by key name
// Accepts the bucket UID and key names for the desired bucket items
// Returns an error
func (b *BucketItemService) DeleteBucketItemsByKeyName(bucketUID string, keys []string) error {
	if utils.IsStringEmpty(bucketUID) {
		return ErrBucketUIDIsEmpty
	}
	for _, key := range keys {
		err := b.bucketItemRepo.DeleteBucketItemByKeyName(bucketUID, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletes all the bucket items for a specific bucket
// Accepts the bucket UID of the bucket
// Returns an error
func (b *BucketItemService) DeleteBucketItems(bucketUID string) error {
	if utils.IsStringEmpty(bucketUID) {
		return ErrBucketUIDIsEmpty
	}
	err := b.bucketItemRepo.DeleteBucketItems(bucketUID)
	if err != nil {
		return err
	}
	return nil
}
