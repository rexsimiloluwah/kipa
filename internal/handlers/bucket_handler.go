package handlers

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/services"
	"keeper/internal/validators"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type BucketHandler struct {
	bucketSvc services.IBucketService
	validator validators.IValidator
}

type IBucketHandler interface {
	CreateBucket(c echo.Context) error
	FindBucketByUID(c echo.Context) error
	ListUserBuckets(c echo.Context) error
	UpdateBucket(c echo.Context) error
	DeleteBucket(c echo.Context) error
}

func NewBucketHandler(cfg *config.Config, dbClient *mongo.Client) IBucketHandler {
	bucketRepo := repository.NewBucketRepository(cfg, dbClient)
	bucketItemRepo := repository.NewBucketItemRepository(cfg, dbClient)
	bucketService := services.NewBucketService(cfg, bucketRepo, bucketItemRepo)
	return &BucketHandler{
		bucketSvc: bucketService,
		validator: validators.NewValidator(),
	}
}

func (h *BucketHandler) CreateBucket(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	q := c.QueryParam("full")
	data := new(dto.CreateBucketInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	resp, err := h.bucketSvc.CreateBucket(*data, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// if 'full' is true, return the full response
	if full, _ := strconv.ParseBool(q); full {
		return c.JSON(http.StatusCreated, &models.SuccessResponse{
			Status:  true,
			Message: "Successfully created a new bucket!",
			Data:    resp,
		})
	}
	return c.JSON(http.StatusCreated, resp.UID)
}

func (h *BucketHandler) FindBucketByUID(c echo.Context) error {
	// retrieve the bucketUID
	bucketUID := c.Param("bucketUID")
	bucket, err := h.bucketSvc.FindBucketByUID(bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched bucket!",
		Data:    bucket,
	})
}

func (h *BucketHandler) ListUserBuckets(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	buckets, err := h.bucketSvc.ListUserBuckets(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched user's buckets!",
		Data:    buckets,
	})
}

func (h *BucketHandler) UpdateBucket(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	data := new(dto.UpdateBucketInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	err := h.bucketSvc.UpdateBucket(bucketUID, *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully updated bucket!"})
}

func (h *BucketHandler) DeleteBucket(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	err := h.bucketSvc.DeleteBucket(bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully deleted bucket!"})
}
