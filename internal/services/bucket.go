package services

import (
	"keeper/internal/config"
	"keeper/internal/repository"
)

type BucketService struct {
	bucketRepository repository.BucketRepository
	cfg              *config.Config
}

type IBucketService interface {
}
