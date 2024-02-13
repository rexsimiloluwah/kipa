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

	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BucketService struct {
	bucketRepo     repository.IBucketRepository
	bucketItemRepo repository.IBucketItemRepository
	cfg            *config.Config
}

type IBucketService interface {
	CreateBucket(data dto.CreateBucketInputDTO, userID primitive.ObjectID) (*dto.CreateBucketOutputDTO, error)
	FindBucketByID(id string) (*dto.BucketDetailsOutput, error)
	FindBucketByUID(uid string) (*dto.BucketDetailsOutput, error)
	ListUserBuckets(userID string) ([]dto.BucketDetailsOutput, error)
	ListUserBucketsPaged(userID string, queryParams url.Values) ([]dto.BucketDetailsOutput, utils.PageInfo, error)
	UpdateBucket(uid string, data dto.UpdateBucketInputDTO) error
	DeleteBucket(uid string) error
}

func NewBucketService(cfg *config.Config, bucketRepo repository.IBucketRepository, bucketItemRepo repository.IBucketItemRepository) IBucketService {
	return &BucketService{
		bucketRepo:     bucketRepo,
		bucketItemRepo: bucketItemRepo,
		cfg:            cfg,
	}
}

var (
	ErrBucketNameIsEmpty = errors.New("bucket name cannot be empty")
	ErrBucketUIDIsEmpty  = errors.New("bucket uid cannot be empty")
	ErrBucketIDIsEmpty   = errors.New("bucket id cannot be empty")
	ErrUserIDIsEmpty     = errors.New("user id cannot be empty")
)

// Service for creating a new bucket
func (b *BucketService) CreateBucket(data dto.CreateBucketInputDTO, userID primitive.ObjectID) (*dto.CreateBucketOutputDTO, error) {
	if utils.IsStringEmpty(data.Name) {
		return &dto.CreateBucketOutputDTO{}, ErrBucketNameIsEmpty
	}

	var permissions models.BucketPermissionsList
	if len(data.Permissions) == 0 {
		// set the default list of bucket permissions
		permissions = models.BUCKET_PERMISSIONS
	} else {
		permissions = data.Permissions
	}

	newBucket := &models.Bucket{
		Name:        data.Name,
		Description: data.Description,
		Permissions: permissions,
		UserID:      userID,
	}
	// Generate a new bucket UID
	newBucketUID := uniuri.NewLen(16)
	newBucket.UID = newBucketUID
	newBucket.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newBucket.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	id, err := b.bucketRepo.CreateBucket(newBucket)
	if err != nil {
		logrus.WithError(err).Error("error saving bucket to database")
		return &dto.CreateBucketOutputDTO{}, fmt.Errorf("error saving bucket to database: %s", err.Error())
	}

	return &dto.CreateBucketOutputDTO{
		ID:          id,
		UID:         newBucketUID,
		Name:        newBucket.Name,
		Description: newBucket.Description,
		Permissions: newBucket.Permissions,
		CreatedAt:   newBucket.CreatedAt,
	}, nil
}

// Service for returning a bucket's details by ID
func (b *BucketService) FindBucketByID(id string) (*dto.BucketDetailsOutput, error) {
	if utils.IsStringEmpty(id) {
		return &dto.BucketDetailsOutput{}, ErrBucketIDIsEmpty
	}
	bucket, err := b.bucketRepo.FindBucketByID(id)
	if err != nil {
		return &dto.BucketDetailsOutput{}, err
	}
	// find the bucket items
	bucketItems, err := b.bucketItemRepo.FindBucketItems(bucket.UID)
	if err != nil && !errors.Is(err, models.ErrBucketItemNotFound) {
		return &dto.BucketDetailsOutput{}, err
	}
	return &dto.BucketDetailsOutput{
		ID:          bucket.ID,
		UID:         bucket.UID,
		UserID:      bucket.UserID,
		Name:        bucket.Name,
		Description: bucket.Description,
		Permissions: bucket.Permissions,
		CreatedAt:   bucket.CreatedAt,
		UpdatedAt:   bucket.UpdatedAt,
		BucketItems: bucketItems,
	}, nil
}

// Service for returning a bucket's details by UID
func (b *BucketService) FindBucketByUID(uid string) (*dto.BucketDetailsOutput, error) {
	if utils.IsStringEmpty(uid) {
		return &dto.BucketDetailsOutput{}, ErrBucketUIDIsEmpty
	}
	bucket, err := b.bucketRepo.FindBucketByUID(uid)
	if err != nil {
		return &dto.BucketDetailsOutput{}, err
	}
	// find the bucket items
	bucketItems, err := b.bucketItemRepo.FindBucketItems(bucket.UID)
	if err != nil && !errors.Is(err, models.ErrBucketItemsNotFound) {
		return &dto.BucketDetailsOutput{}, err
	}
	return &dto.BucketDetailsOutput{
		ID:          bucket.ID,
		UID:         bucket.UID,
		UserID:      bucket.UserID,
		Name:        bucket.Name,
		Description: bucket.Description,
		Permissions: bucket.Permissions,
		CreatedAt:   bucket.CreatedAt,
		UpdatedAt:   bucket.UpdatedAt,
		BucketItems: bucketItems,
	}, nil
}

// Service for listing all a user's buckets with bucket items (using pagination and filtering)
func (b *BucketService) ListUserBucketsPaged(userID string, queryParams url.Values) ([]dto.BucketDetailsOutput, utils.PageInfo, error) {
	if utils.IsStringEmpty(userID) {
		return nil, utils.PageInfo{}, ErrUserIDIsEmpty
	}
	userBucketDetailsOutput := []dto.BucketDetailsOutput{}
	filter, findOpts, paginationParams, err := utils.ParseRequestQueryParams(queryParams)
	if err != nil {
		return nil, utils.PageInfo{}, err
	}
	// find all user's buckets
	userBuckets, pageInfo, err := b.bucketRepo.FindBucketsByUserIDPaged(userID, filter, findOpts, paginationParams)
	if err != nil {
		return nil, utils.PageInfo{}, err
	}
	// construct the bucket detals response
	for _, bucket := range userBuckets {
		bucketItems, err := b.bucketItemRepo.FindBucketItems(bucket.UID)
		if err != nil && !errors.Is(err, models.ErrBucketItemsNotFound) {
			logrus.WithError(err).Errorf("error fetching bucket items for bucket: %s", bucket.UID)
			return nil, utils.PageInfo{}, err
		}
		bucketDetails := dto.BucketDetailsOutput{
			ID:          bucket.ID,
			UID:         bucket.UID,
			UserID:      bucket.UserID,
			Name:        bucket.Name,
			Description: bucket.Description,
			Permissions: bucket.Permissions,
			CreatedAt:   bucket.CreatedAt,
			UpdatedAt:   bucket.UpdatedAt,
			BucketItems: bucketItems,
		}
		// append single bucket's details to the final response
		userBucketDetailsOutput = append(userBucketDetailsOutput, bucketDetails)
	}
	return userBucketDetailsOutput, pageInfo, nil
}

// Service for listing all a user's buckets with bucket items
func (b *BucketService) ListUserBuckets(userID string) ([]dto.BucketDetailsOutput, error) {
	if utils.IsStringEmpty(userID) {
		return nil, ErrUserIDIsEmpty
	}
	userBucketDetailsOutput := []dto.BucketDetailsOutput{}
	// find all user's buckets
	userBuckets, err := b.bucketRepo.FindBucketsByUserID(userID)
	if err != nil {
		return nil, err
	}
	// construct the bucket detals response
	for _, bucket := range userBuckets {
		bucketItems, err := b.bucketItemRepo.FindBucketItems(bucket.UID)
		if err != nil && !errors.Is(err, models.ErrBucketItemsNotFound) {
			logrus.WithError(err).Errorf("error fetching bucket items for bucket: %s", bucket.UID)
			return nil, err
		}
		bucketDetails := dto.BucketDetailsOutput{
			ID:          bucket.ID,
			UID:         bucket.UID,
			UserID:      bucket.UserID,
			Name:        bucket.Name,
			Description: bucket.Description,
			Permissions: bucket.Permissions,
			CreatedAt:   bucket.CreatedAt,
			UpdatedAt:   bucket.UpdatedAt,
			BucketItems: bucketItems,
		}
		// append single bucket's details to the final response
		userBucketDetailsOutput = append(userBucketDetailsOutput, bucketDetails)
	}
	return userBucketDetailsOutput, nil
}

func (b *BucketService) UpdateBucket(uid string, data dto.UpdateBucketInputDTO) error {
	if utils.IsStringEmpty(uid) {
		return ErrBucketUIDIsEmpty
	}
	updatedBucket := &models.Bucket{
		Name:        data.Name,
		Description: data.Description,
		Permissions: data.Permissions,
		UID:         uid,
	}
	updatedBucket.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	err := b.bucketRepo.UpdateBucket(updatedBucket)
	if err != nil {
		return err
	}
	return nil
}

func (b *BucketService) DeleteBucket(uid string) error {
	if utils.IsStringEmpty(uid) {
		return ErrBucketUIDIsEmpty
	}
	// delete all the bucket items
	err := b.bucketItemRepo.DeleteBucketItems(uid)
	if err != nil {
		return err
	}
	err = b.bucketRepo.DeleteBucketByUID(uid)
	if err != nil {
		return err
	}
	return nil
}
