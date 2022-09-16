package handlers

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/services"
	"keeper/internal/validators"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIKeyHandler struct {
	apiKeySvc services.IAPIKeyService
	validator validators.IValidator
}

type IAPIKeyHandler interface {
	CreateAPIKey(c echo.Context) error
	FindAPIKeyByID(c echo.Context) error
	FindUserAPIKeys(c echo.Context) error
	UpdateAPIKey(c echo.Context) error
	RevokeAPIKey(c echo.Context) error
	DeleteAPIKey(c echo.Context) error
	RevokeAPIKeys(c echo.Context) error
	DeleteAPIKeys(c echo.Context) error
}

func NewAPIKeyHandler(cfg *config.Config, dbClient *mongo.Client) IAPIKeyHandler {
	apiKeyRepo := repository.NewAPIKeyRepository(cfg, dbClient)
	apiKeyService := services.NewAPIKeyService(cfg, apiKeyRepo)
	return &APIKeyHandler{
		apiKeySvc: apiKeyService,
		validator: validators.NewValidator(),
	}
}

// CreateAPIKey godoc
// @Summary      CreateAPIKey
// @Description  Create a new API key
// @Tags         APIKey
// @Produce      json
// @Param data body dto.CreateAPIKeyInputDTO true "Create API Key Data"
// @Security BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /api_key [post]
func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	data := new(dto.CreateAPIKeyInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	resp, err := h.apiKeySvc.CreateAPIKey(*data, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully created a new API Key!",
		Data:    resp,
	})
}

func (h *APIKeyHandler) FindAPIKeyByID(c echo.Context) error {
	// retrieve the apiKeyID
	apiKeyID := c.Param("apiKeyId")
	apiKey, err := h.apiKeySvc.FindAPIKeyByID(apiKeyID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully found API Key!",
		Data:    apiKey,
	})
}

func (h *APIKeyHandler) FindUserAPIKeys(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	apiKeys, err := h.apiKeySvc.FindUserAPIKeys(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched user's API Keys!",
		Data:    apiKeys,
	})
}

func (h *APIKeyHandler) UpdateAPIKey(c echo.Context) error {
	// retrieve apiKeyID from param
	apiKeyId := c.Param("apiKeyId")
	data := new(dto.UpdateAPIKeyInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	err := h.apiKeySvc.UpdateAPIKey(apiKeyId, *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully updated API key!",
	})
}

func (h *APIKeyHandler) RevokeAPIKeys(c echo.Context) error {
	apiKeyIds := new(dto.APIKeysIDsInputDTO)
	if err := c.Bind(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	err := h.apiKeySvc.RevokeAPIKeys(apiKeyIds.Ids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully revoked API Key(s)!"})
}

func (h *APIKeyHandler) DeleteAPIKeys(c echo.Context) error {
	apiKeyIds := new(dto.APIKeysIDsInputDTO)
	if err := c.Bind(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	err := h.apiKeySvc.DeleteAPIKeys(apiKeyIds.Ids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully deleted API Key(s)!"})
}

func (h *APIKeyHandler) RevokeAPIKey(c echo.Context) error {
	// retrieve apiKeyID from param
	apiKeyId := c.Param("apiKeyId")
	err := h.apiKeySvc.RevokeAPIKey(apiKeyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully revoked API key!",
	})
}

func (h *APIKeyHandler) DeleteAPIKey(c echo.Context) error {
	// retrieve apiKeyID from param
	apiKeyId := c.Param("apiKeyId")
	err := h.apiKeySvc.DeleteAPIKey(apiKeyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully deleted API key!",
	})
}
