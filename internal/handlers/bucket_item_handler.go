package handlers

import (
	"fmt"
	"io/ioutil"
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

type BucketItemHandler struct {
	bucketItemSvc services.IBucketItemService
	validator     validators.IValidator
}

type IBucketItemHandler interface {
	CreateBucketItem(c echo.Context) error
	FindBucketItemByID(c echo.Context) error
	ListBucketItems(c echo.Context) error
	FindBucketItemByKeyName(c echo.Context) error
	UpdateBucketItemByKeyName(c echo.Context) error
	DeleteBucketItemById(c echo.Context) error
	DeleteBucketItemsById(c echo.Context) error
	DeleteBucketItemByKeyName(c echo.Context) error
}

func NewBucketItemHandler(cfg *config.Config, dbClient *mongo.Client) IBucketItemHandler {
	bucketRepo := repository.NewBucketRepository(cfg, dbClient)
	bucketItemRepo := repository.NewBucketItemRepository(cfg, dbClient)
	bucketItemService := services.NewBucketItemService(cfg, bucketItemRepo, bucketRepo)
	return &BucketItemHandler{
		bucketItemSvc: bucketItemService,
	}
}

func (h *BucketItemHandler) CreateBucketItem(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	q := c.QueryParam("full")
	// retrieve user from context
	user := c.Get("user").(*models.User)
	data := new(dto.CreateBucketItemInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	resp, err := h.bucketItemSvc.CreateBucketItem(*data, user.ID, bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	if full, _ := strconv.ParseBool(q); full {
		return c.JSON(http.StatusCreated, &models.SuccessResponse{
			Status:  true,
			Message: fmt.Sprintf("Successfully created '%s' in bucket '%s'", data.Key, bucketUID),
			Data:    resp,
		})
	}
	return c.JSON(http.StatusCreated, fmt.Sprintf("Successfully created '%s' in bucket '%s'", data.Key, bucketUID))
}

func (h *BucketItemHandler) FindBucketItemByID(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (h *BucketItemHandler) ListBucketItems(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	bucketItems, err := h.bucketItemSvc.ListBucketItems(bucketUID)
	fmt.Println(bucketItems)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched bucket items!",
		Data:    bucketItems,
	})
}

func (h *BucketItemHandler) FindBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	q := c.QueryParam("full")
	// retrieve the key
	key := c.Param("key")
	bucketItem, err := h.bucketItemSvc.FindBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	if full, _ := strconv.ParseBool(q); full {
		return c.JSON(http.StatusOK, &models.SuccessResponse{
			Status:  true,
			Message: fmt.Sprintf("Successfully fetched '%s'!", key),
			Data:    bucketItem,
		})
	}
	return c.JSON(http.StatusOK, bucketItem.Data)
}

func (h *BucketItemHandler) UpdateBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	// retrieve the key
	key := c.Param("key")
	data := new(dto.UpdateBucketItemInputDTO)
	if err := c.Bind(data); err != nil {
		// handle for increments/decrements
		body, _ := ioutil.ReadAll(c.Request().Body)
		value := string(body)
		amount, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
		}
		err = h.bucketItemSvc.IncrementIntValue(bucketUID, key, int(amount))
		if err != nil {
			return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
		}
		return c.JSON(http.StatusOK, &models.SuccessResponse{
			Status:  true,
			Message: fmt.Sprintf("Successfully updated '%s'!", key),
		})
	}
	err := h.bucketItemSvc.UpdateBucketItemByKeyName(*data, bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: fmt.Sprintf("Successfully updated '%s'!", key),
	})
}

func (h *BucketItemHandler) DeleteBucketItemById(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (h *BucketItemHandler) DeleteBucketItemsById(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (h *BucketItemHandler) DeleteBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	// retrieve the key
	key := c.Param("key")
	err := h.bucketItemSvc.DeleteBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: fmt.Sprintf("Successfully deleted '%s'", key),
	})
}
