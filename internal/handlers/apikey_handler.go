package handlers

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIKeyHandler struct {
	APIKeySvc services.IAPIKeyService
}

type IAPIKeyHandler interface {
	CreateAPIKey(c echo.Context) error
	FindAPIKeyByID(c echo.Context) error
	FindUserAPIKeys(c echo.Context) error
	UpdateAPIKey(c echo.Context) error
	RevokeAPIKeys(c echo.Context) error
	DeleteAPIKeys(c echo.Context) error
}

func NewAPIKeyHandler(dbClient *mongo.Client) IAPIKeyHandler {
	cfg := config.New()
	apiKeyService := services.NewAPIKeyService(cfg, dbClient)
	return &APIKeyHandler{
		APIKeySvc: apiKeyService,
	}
}

func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	data := new(dto.CreateAPIKeyInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	resp, err := h.APIKeySvc.CreateAPIKey(*data, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"status": true, "message": "Successfully created a new API Key!", "data": resp})
}

func (h *APIKeyHandler) FindAPIKeyByID(c echo.Context) error {
	// retrieve the apiKeyID
	apiKeyID := c.Param("apiKeyId")
	apiKey, err := h.APIKeySvc.FindAPIKeyByID(apiKeyID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully found API Key!", "data": apiKey})
}

func (h *APIKeyHandler) FindUserAPIKeys(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)
	apiKeys, err := h.APIKeySvc.FindUserAPIKeys(user.ID.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched user's API Keys!", "data": apiKeys})
}

func (h *APIKeyHandler) UpdateAPIKey(c echo.Context) error {
	// retrieve apiKeyID from param
	apiKeyId := c.Param("apiKeyId")
	data := new(dto.UpdateAPIKeyInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	err := h.APIKeySvc.UpdateAPIKey(apiKeyId, *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully updated API key!"})
}

func (h *APIKeyHandler) RevokeAPIKeys(c echo.Context) error {
	apiKeyIds := make([]string, 0)
	if err := c.Bind(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	err := h.APIKeySvc.RevokeAPIKeys(apiKeyIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully revoked API Key(s)!"})
}

func (h *APIKeyHandler) DeleteAPIKeys(c echo.Context) error {
	apiKeyIds := make([]string, 0)
	if err := c.Bind(apiKeyIds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	err := h.APIKeySvc.DeleteAPIKeys(apiKeyIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully deleted API Key(s)!"})
}
