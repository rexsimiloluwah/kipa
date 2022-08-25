package handlers

import (
	"fmt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type BucketItemHandler struct {
	BucketItemSvc services.IBucketItemService
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

func NewBucketItemHandler(dbClient *mongo.Client) IBucketItemHandler {
	cfg := config.New()
	bucketRepo := repository.NewBucketRepository(cfg, dbClient)
	bucketItemRepo := repository.NewBucketItemRepository(cfg, dbClient)
	bucketItemService := services.NewBucketItemService(cfg, bucketItemRepo, bucketRepo)
	return &BucketItemHandler{
		BucketItemSvc: bucketItemService,
	}
}

func (h *BucketItemHandler) CreateBucketItem(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	// retrieve user from context
	user := c.Get("user").(*models.User)
	data := new(dto.CreateBucketItemInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	resp, err := h.BucketItemSvc.CreateBucketItem(*data, user.ID, bucketUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"status": true, "message": "Successfully created new item!", "data": resp})
}

func (h *BucketItemHandler) FindBucketItemByID(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (h *BucketItemHandler) ListBucketItems(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (h *BucketItemHandler) FindBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	// retrieve the key
	key := c.Param("key")
	bucketItem, err := h.BucketItemSvc.FindBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched item!", "data": bucketItem})
}

func (h *BucketItemHandler) UpdateBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	// retrieve the key
	key := c.Param("key")
	data := new(dto.UpdateBucketItemInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	err := h.BucketItemSvc.UpdateBucketItemByKeyName(*data, bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully updated item!"})
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
	err := h.BucketItemSvc.DeleteBucketItemByKeyName(bucketUID, key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": fmt.Sprintf("Successfully deleted %s", key)})
}
