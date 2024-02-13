package handlers

import (
	"fmt"
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
	ListUserBucketsPaged(c echo.Context) error
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

// CreateBucket  godoc
// @Summary      CreateBucket
// @Description  Create a new bucket
// @Tags         Bucket
// @Produce      json
// @Param        data body dto.CreateBucketInputDTO true "Create Bucket Data"
// @Param        full query bool false "Should return full response"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /bucket [post]
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

// FindBucketByUID  godoc
// @Summary      FindBucketByUID
// @Description  Returns a bucket that matches the passed UID
// @Tags         Bucket
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /bucket/{bucketUID} [get]
func (h *BucketHandler) FindBucketByUID(c echo.Context) error {
	// retrieve the bucketUID
	bucketUID := c.Param("bucketUID")
	bucket, err := h.bucketSvc.FindBucketByUID(bucketUID)
	if err != nil {
		if err == models.ErrBucketNotFound {
			return c.JSON(http.StatusNotFound, &models.ErrorResponse{Status: false, Error: err.Error()})
		}
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched bucket!",
		Data:    bucket,
	})
}

// ListUserBuckets  godoc
// @Summary      ListUserBuckets
// @Description  Returns a list of the authenticated user's buckets
// @Tags         Bucket
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /buckets/all [get]
func (h *BucketHandler) ListUserBuckets(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	buckets, err := h.bucketSvc.ListUserBuckets(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: fmt.Sprintf("Successfully fetched %d buckets!", len(buckets)),
		Data:    buckets,
	})
}

// ListUserBucketsPaged  godoc
// @Summary      ListUserBucketsPaged
// @Description  Returns a list of the authenticated user's buckets (with pagination, sorting, and filtering)
// @Tags         Bucket
// @Produce      json
// @Param        page query integer false "Current Page"
// @Param        perPage query integer false "Per Page"
// @Param        sortBy query string false "Sort By"
// @Security     BearerAuth
// @Success      200  {object} 	models.PaginatedSuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /buckets [get]
func (h *BucketHandler) ListUserBucketsPaged(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	buckets, pageInfo, err := h.bucketSvc.ListUserBucketsPaged(user.ID.Hex(), c.QueryParams())
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.PaginatedSuccessResponse{
		Status:   true,
		Message:  fmt.Sprintf("Successfully fetched %d buckets!", len(buckets)),
		Data:     buckets,
		PageInfo: pageInfo,
	})
}

// UpdateBucket  godoc
// @Summary      UpdateBucket
// @Description  Update a user's bucket
// @Tags         Bucket
// @Produce      json
// @Param        data body dto.UpdateBucketInputDTO true "Update Bucket Data"
// @Param        bucketUID path string true "Bucket UID"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /bucket/{bucketUID} [put]
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

// DeleteBucket  godoc
// @Summary      DeleteBucket
// @Description  Delete a user's bucket
// @Tags         Bucket
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /bucket/{bucketUID} [delete]
func (h *BucketHandler) DeleteBucket(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	err := h.bucketSvc.DeleteBucket(bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully deleted bucket!"})
}
