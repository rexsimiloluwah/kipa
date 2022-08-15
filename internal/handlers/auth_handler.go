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

type AuthHandler struct {
	AuthSvc services.IAuthService
}

type IAuthHandler interface {
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
	GetAuthUser(c echo.Context) error
}

func NewAuthHandler(dbClient *mongo.Client) IAuthHandler {
	cfg := config.New()
	authService := services.NewAuthService(cfg, dbClient)
	return &AuthHandler{
		AuthSvc: authService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	data := new(dto.LoginUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}

	resp, err := h.AuthSvc.Login(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully logged in user!", "data": resp})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)

	resp, err := h.AuthSvc.RefreshToken(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully refreshed token!", "data": resp})
}

func (h *AuthHandler) GetAuthUser(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched authenticated user.", "data": user})
}
