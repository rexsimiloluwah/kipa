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
	ListBucketItemsByBucketUID(c echo.Context) error
	ListBucketItemsPaged(c echo.Context) error
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

// CreateBucketItem  godoc
// @Summary      CreateBucketItem
// @Description  Create a new bucket item
// @Tags         BucketItem
// @Produce      json
// @Param        data body dto.CreateBucketItemInputDTO true "Create Bucket Item Data"
// @Param        bucketUID path string true "Bucket UID"
// @Param        full query bool false "Should return full response"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /item/{bucketUID} [post]
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
	// if err := h.validator.Validate(data); err != nil {
	// 	return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	// }
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

// ListBucketItems  godoc
// @Summary      ListBucketItems
// @Description  Returns a list of all the items contained in a bucket
// @Tags         BucketItem
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Param        full query bool false "Should return full response"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /items/{bucketUID} [get]
func (h *BucketItemHandler) ListBucketItemsByBucketUID(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	bucketItems, err := h.bucketItemSvc.ListBucketItems(bucketUID)

	fmt.Println(bucketItems)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: fmt.Sprintf("Successfully fetched %d bucket items!", len(bucketItems)),
		Data:    bucketItems,
	})
}

// ListBucketItemsPaged  godoc
// @Summary      ListBucketItemsPaged
// @Description  Returns a list of all the items contained in a bucket (supports pagination, filtering, and sorting)
// @Tags         BucketItem
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Param        full query bool false "Should return full response"
// @Param        page query integer false "Current Page"
// @Param        perPage query integer false "Per Page"
// @Param        sortBy query string false "Sort By"
// @Security     BearerAuth
// @Success      200  {object} 	models.PaginatedSuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /items/{bucketUID} [get]
func (h *BucketItemHandler) ListBucketItemsPaged(c echo.Context) error {
	queryParams := c.QueryParams()
	bucketItems, pageInfo, err := h.bucketItemSvc.ListBucketItemsPaged(queryParams)

	fmt.Println(bucketItems)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.PaginatedSuccessResponse{
		Status:   true,
		Message:  fmt.Sprintf("Successfully fetched %d bucket items!", len(bucketItems)),
		Data:     bucketItems,
		PageInfo: pageInfo,
	})
}

// FindBucketItemByKeyName  godoc
// @Summary      FindBucketItemByKeyName
// @Description  Returns an item from a bucket matching the passed key
// @Tags         BucketItem
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Param        key path string true "Key name"
// @Param        full query bool false "Should return full response"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /item/{bucketUID}/{key} [get]
func (h *BucketItemHandler) FindBucketItemByKeyName(c echo.Context) error {
	// retrieve the bucket UID
	bucketUID := c.Param("bucketUID")
	q := c.QueryParam("full")
	// retrieve the key
	key := c.Param("key")
	bucketItem, err := h.bucketItemSvc.FindBucketItemByKeyName(bucketUID, key)
	if err != nil {
		if err == models.ErrBucketItemNotFound {
			return c.JSON(http.StatusNotFound, &models.ErrorResponse{Status: false, Error: err.Error()})
		}
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

// UpdateBucketItemByKeyName  godoc
// @Summary      UpdateBucketItemByKeyName
// @Description  Update an item from a bucket matching the passed key
// @Tags         BucketItem
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Param        key path string true "Key name"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /items/{bucketUID}/{key} [put]
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

// DeleteBucketItemByKeyName  godoc
// @Summary      DeleteBucketItemByKeyName
// @Description  Delete an item from a bucket matching the passed key
// @Tags         BucketItem
// @Produce      json
// @Param        bucketUID path string true "Bucket UID"
// @Param        key path string true "Key name"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /items/{bucketUID}/{key} [delete]
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
