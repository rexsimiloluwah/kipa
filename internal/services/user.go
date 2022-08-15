package services

import (
	"errors"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	UserRepository repository.IUserRepository
	Cfg            *config.Config
}

type IUserService interface {
	Register(data dto.CreateUserInputDTO) error
	FindUserByID(id string) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	UpdateUser(id string, data dto.UpdateUserInputDTO) error
	DeleteUser(id string) error
}

func NewUserService(cfg *config.Config, dbClient *mongo.Client) IUserService {
	userRepository := repository.NewUserRepository(cfg, dbClient)
	return &UserService{
		UserRepository: userRepository,
		Cfg:            cfg,
	}
}

func (s *UserService) Register(data dto.CreateUserInputDTO) error {
	existingUser, err := s.UserRepository.FindUserByEmail(data.Email)
	// check if the user already exists
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return err
	}
	if existingUser != nil {
		return models.ErrUserAlreadyExists
	}
	newUser := &models.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Username:  data.Username,
		Email:     data.Email,
	}

	newUser.ID = primitive.NewObjectID()
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}
	newUser.Password = hashedPassword
	newUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// save the user in the database
	if err := s.UserRepository.CreateUser(newUser); err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindUserByID(id string) (*models.User, error) {
	user, err := s.UserRepository.FindUserById(id)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (s *UserService) FindAllUsers() ([]models.User, error) {
	users, err := s.UserRepository.FindAllUsers()
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(id string, data dto.UpdateUserInputDTO) error {
	user := &models.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Username:  data.Username,
	}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidObjectID
	}
	user.ID = ID
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	if err := s.UserRepository.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id string) error {
	if err := s.UserRepository.DeleteUser(id); err != nil {
		return err
	}
	return nil
}
