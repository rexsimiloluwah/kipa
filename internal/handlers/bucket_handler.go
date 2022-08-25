package handlers

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type BucketHandler struct {
	BucketSvc services.IBucketService
}

type IBucketHandler interface {
	CreateBucket(c echo.Context) error
	FindBucketByUID(c echo.Context) error
	ListUserBuckets(c echo.Context) error
	UpdateBucket(c echo.Context) error
	DeleteBucket(c echo.Context) error
}

func NewBucketHandler(dbClient *mongo.Client) IBucketHandler {
	cfg := config.New()
	bucketRepo := repository.NewBucketRepository(cfg, dbClient)
	bucketItemRepo := repository.NewBucketItemRepository(cfg, dbClient)
	bucketService := services.NewBucketService(cfg, bucketRepo, bucketItemRepo)
	return &BucketHandler{
		BucketSvc: bucketService,
	}
}

func (h *BucketHandler) CreateBucket(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	data := new(dto.CreateBucketInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	resp, err := h.BucketSvc.CreateBucket(*data, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"status": true, "message": "Successfully created a new bucket!", "data": resp})
}

func (h *BucketHandler) FindBucketByUID(c echo.Context) error {
	// retrieve the bucketUID
	bucketUID := c.Param("bucketUID")
	bucket, err := h.BucketSvc.FindBucketByUID(bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched bucket!", "data": bucket})
}

func (h *BucketHandler) ListUserBuckets(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	buckets, err := h.BucketSvc.ListUserBuckets(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched user's buckets!", "data": buckets})
}

func (h *BucketHandler) UpdateBucket(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	data := new(dto.UpdateBucketInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	err := h.BucketSvc.UpdateBucket(bucketUID, *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully updated bucket!"})
}

func (h *BucketHandler) DeleteBucket(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	err := h.BucketSvc.DeleteBucket(bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully deleted bucket!"})
}
