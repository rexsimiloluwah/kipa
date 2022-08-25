package handlers

import (
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserHandler interface {
	GetUserByID(c echo.Context) error
	GetAllUsers(c echo.Context) error
	RegisterUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type UserHandler struct {
	userSvc services.IUserService
}

func NewUserHandler(dbClient *mongo.Client) IUserHandler {
	cfg := config.New()
	userRepo := repository.NewUserRepository(cfg, dbClient)
	userService := services.NewUserService(cfg, userRepo)
	return &UserHandler{
		userSvc: userService,
	}
}

func (u *UserHandler) GetUserByID(c echo.Context) error {
	userId := c.Param("userId")
	user, err := u.userSvc.FindUserByID(userId)
	if err != nil {
		return c.JSON(400, map[string]interface{}{"status": false, "error": err.Error()})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"status": false, "error": "user not found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched user!", "data": user})
}

func (u *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := u.userSvc.FindAllUsers()
	if err != nil {
		return c.JSON(400, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully fetched all users!", "data": users})
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	data := new(dto.CreateUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	err := u.userSvc.Register(*data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"status": true, "message": "Successfully registered user!"})
}

func (u *UserHandler) UpdateUser(c echo.Context) error {
	data := new(dto.UpdateUserInputDTO)
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	// retrieve user from context
	user := c.Get("user").(*models.User)
	err := u.userSvc.UpdateUser(user.ID.Hex(), *data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully updated user!"})
}

func (u *UserHandler) DeleteUser(c echo.Context) error {
	// retrieve user from context
	user := c.Get("user").(*models.User)
	err := u.userSvc.DeleteUser(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Successfully deleted user!"})
}
