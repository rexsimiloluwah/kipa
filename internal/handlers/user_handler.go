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

type IUserHandler interface {
	GetUserByID(c echo.Context) error
	GetAllUsers(c echo.Context) error
	UpdateUser(c echo.Context) error
	UpdateUserPassword(c echo.Context) error
	DeleteUser(c echo.Context) error
	VerifyEmail(c echo.Context) error
}

type UserHandler struct {
	userSvc   services.IUserService
	validator validators.IValidator
}

func NewUserHandler(cfg *config.Config, dbClient *mongo.Client) IUserHandler {
	userRepo := repository.NewUserRepository(cfg, dbClient)
	userService := services.NewUserService(cfg, userRepo)
	return &UserHandler{
		userSvc:   userService,
		validator: validators.NewValidator(),
	}
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	userId := c.Param("userId")
	user, err := h.userSvc.FindUserByID(userId)
	if err != nil {
		return c.JSON(400, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, &models.ErrorResponse{Status: false, Error: "user not found"})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched user!",
		Data:    user,
	})
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.userSvc.FindAllUsers()
	if err != nil {
		return c.JSON(400, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully fetched all users!",
		Data:    users,
	})
}

// VerifyEmail godoc
// @Summary      Verify a User's Email
// @Description  Verify a User's Email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.VerifyEmailInputDTO true "Email Verification Data"
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /user/verify-email [post]
func (h *UserHandler) VerifyEmail(c echo.Context) error {
	data := new(dto.VerifyEmailInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}

	err := h.userSvc.VerifyEmail(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Email verification successful!"})
}

// UpdateUser godoc
// @Summary      UpdateUser
// @Description  Update data for a user
// @Tags         User
// @Produce      json
// @Param        data body dto.UpdateUserInputDTO true "Update User Data"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /user [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	data := new(dto.UpdateUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// retrieve user from context
	user := c.Get("user").(*models.User)
	err := h.userSvc.UpdateUser(user.ID.Hex(), *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{Status: true, Message: "Successfully updated user!"})
}

// UpdateUserPassword godoc
// @Summary      UpdateUserPassword
// @Description  Update a user's password
// @Tags         User
// @Produce      json
// @Param        data body dto.UpdateUserPasswordInputDTO true "Update User Password Data"
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /user/password [put]
func (h *UserHandler) UpdateUserPassword(c echo.Context) error {
	data := new(dto.UpdateUserPasswordInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// validate the request data
	if err := h.validator.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	// retrieve user from context
	user := c.Get("user").(*models.User)
	err := h.userSvc.UpdateUserPassword(user.ID.Hex(), *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully updated user password!",
	})
}

// DeleteUser godoc
// @Summary      DeleteUser
// @Description  Delete a user's account
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} 	models.SuccessResponse
// @Failure      400  {object} 	models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router /user [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	err := h.userSvc.DeleteUser(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.ErrorResponse{Status: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  true,
		Message: "Successfully deleted user!",
	})
}
