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

type AuthHandler struct {
	authSvc   services.IAuthService
	userSvc   services.IUserService
	validator validators.IValidator
}

type IAuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
	GetAuthUser(c echo.Context) error
	ForgotPassword(c echo.Context) error
	ResetPassword(c echo.Context) error
}

func NewAuthHandler(cfg *config.Config, dbClient *mongo.Client) IAuthHandler {
	userRepo := repository.NewUserRepository(cfg, dbClient)
	authService := services.NewAuthService(cfg, userRepo)
	userService := services.NewUserService(cfg, userRepo)
	return &AuthHandler{
		authSvc:   authService,
		userSvc:   userService,
		validator: validators.NewValidator(),
	}
}

// RegisterUser godoc
// @Summary      Register user
// @Description  Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.CreateUserInputDTO true "User Register Data"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	data := new(dto.CreateUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	err := h.userSvc.Register(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, &models.SuccessResponse{Status: true, Message: "Successfully registered user!"})
}

// LoginUser godoc
// @Summary      Login user
// @Description  Login an existing user, returns the access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.LoginUserInputDTO true "User Login Data"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	data := new(dto.LoginUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	resp, err := h.authSvc.Login(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully logged in user!", Data: resp})
}

// RefreshToken godoc
// @Summary      Refresh token
// @Description  Refreshes a user's access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        x-refresh-token header string true "Refresh token"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)

	resp, err := h.authSvc.RefreshToken(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully refreshed token!", Data: resp})
}

// GetAuthUser godoc
// @Summary      Get Auth User
// @Description  Returns the authenticated user decoded from the bearer token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/user [get]
func (h *AuthHandler) GetAuthUser(c echo.Context) error {
	// retrieve the user from context
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully fetched authenticated user.", Data: user})
}

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Sends a reset password link to the user's email if it exists in the database
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.ForgotPasswordInputDTO true "Forgot Password Data"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	data := new(dto.ForgotPasswordInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	err := h.authSvc.ForgotPassword(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully sent reset password link!"})
}

// ResetPassword godoc
// @Summary      Reset Password
// @Description  Reset a user's password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.ResetPasswordInputDTO true "Reset Password Data"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	data := new(dto.ResetPasswordInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	err := h.authSvc.ResetPassword(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully reset password!"})
}
