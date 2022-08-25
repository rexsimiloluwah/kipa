package services

import (
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"time"

	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BucketService struct {
	BucketRepository     repository.IBucketRepository
	BucketItemRepository repository.IBucketItemRepository
	Cfg                  *config.Config
}

type IBucketService interface {
	CreateBucket(data dto.CreateBucketInputDTO, userID primitive.ObjectID) (*dto.CreateBucketOutputDTO, error)
	FindBucketByID(id string) (*dto.BucketDetailsOutput, error)
	FindBucketByUID(uid string) (*dto.BucketDetailsOutput, error)
	ListUserBuckets(userID string) ([]dto.BucketDetailsOutput, error)
	UpdateBucket(uid string, data dto.UpdateBucketInputDTO) error
	DeleteBucket(uid string) error
}

func NewBucketService(cfg *config.Config, bucketRepo repository.IBucketRepository, bucketItemRepo repository.IBucketItemRepository) IBucketService {
	return &BucketService{
		BucketRepository:     bucketRepo,
		BucketItemRepository: bucketItemRepo,
		Cfg:                  cfg,
	}
}

var (
	ErrBucketNameIsEmpty = errors.New("bucket name cannot be empty")
)

// Service for creating a new bucket
func (b *BucketService) CreateBucket(data dto.CreateBucketInputDTO, userID primitive.ObjectID) (*dto.CreateBucketOutputDTO, error) {
	if utils.IsStringEmpty(data.Name) {
		return &dto.CreateBucketOutputDTO{}, ErrBucketNameIsEmpty
	}
	newBucket := &models.Bucket{
		Name:        data.Name,
		Description: data.Description,
		Permissions: data.Permissions,
		UserID:      userID,
	}
	// Generate a new bucket UID
	newBucketUID := uniuri.NewLen(16)
	newBucket.UID = newBucketUID
	newBucket.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newBucket.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	id, err := b.BucketRepository.CreateBucket(newBucket)
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
	bucket, err := b.BucketRepository.FindBucketByID(id)
	if err != nil {
		return &dto.BucketDetailsOutput{}, err
	}
	// find the bucket items
	bucketItems, err := b.BucketItemRepository.FindBucketItems(bucket.UID)
	if err != nil {
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
	bucket, err := b.BucketRepository.FindBucketByUID(uid)
	if err != nil {
		return &dto.BucketDetailsOutput{}, err
	}
	// find the bucket items
	bucketItems, err := b.BucketItemRepository.FindBucketItems(bucket.UID)
	if err != nil {
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

// Service for listing all a user's buckets with bucket items
func (b *BucketService) ListUserBuckets(userID string) ([]dto.BucketDetailsOutput, error) {
	userBucketDetailsOutput := []dto.BucketDetailsOutput{}
	// find all user's buckets
	userBuckets, err := b.BucketRepository.FindBucketsByUserID(userID)
	if err != nil {
		return []dto.BucketDetailsOutput{}, err
	}
	// construct the bucket detals response
	for _, bucket := range userBuckets {
		bucketItems, err := b.BucketItemRepository.FindBucketItems(bucket.UID)
		if err != nil {
			logrus.WithError(err).Errorf("error fetching bucket items for bucket: %s", bucket.UID)
			return []dto.BucketDetailsOutput{}, err
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
	updatedBucket := &models.Bucket{
		Name:        data.Name,
		Description: data.Description,
		Permissions: data.Permissions,
		UID:         uid,
	}
	updatedBucket.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	err := b.BucketRepository.UpdateBucket(updatedBucket)
	if err != nil {
		return err
	}
	return nil
}

func (b *BucketService) DeleteBucket(uid string) error {
	// delete all the bucket items
	err := b.BucketItemRepository.DeleteBucketItems(uid)
	if err != nil {
		return err
	}
	err = b.BucketRepository.DeleteBucketByUID(uid)
	if err != nil {
		return err
	}
	return nil
}
