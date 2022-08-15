package services

import (
	"keeper/internal/config"
	"keeper/internal/repository"
)

type BucketItemService struct {
	bucketItemRepository repository.BucketRepository
	cfg                  *config.Config
}

type IBucketItemService interface {
}
