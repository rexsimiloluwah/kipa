package services

import (
	"errors"
	"fmt"
	"keeper/internal/auth/jwt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/pkg/mailer"
	"keeper/internal/queue"
	"keeper/internal/queue/tasks"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo repository.IUserRepository
	jwtSvc   jwt.IJwtService
	cfg      *config.Config
	queue    *queue.RedisQueue
}

type IUserService interface {
	Register(data dto.CreateUserInputDTO) error
	FindUserByID(id string) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	UpdateUser(id string, data dto.UpdateUserInputDTO) error
	UpdateUserPassword(id string, data dto.UpdateUserPasswordInputDTO) error
	DeleteUser(id string) error
	VerifyEmail(data dto.VerifyEmailInputDTO) error
}

func NewUserService(cfg *config.Config, userRepo repository.IUserRepository) IUserService {
	jwtSvc := jwt.NewJwtService(cfg, userRepo)
	queue := queue.NewRedisQueue(cfg)
	return &UserService{
		userRepo: userRepo,
		cfg:      cfg,
		jwtSvc:   jwtSvc,
		queue:    queue,
	}
}

// Create a new user
func (s *UserService) Register(data dto.CreateUserInputDTO) error {
	if data.Email == "" {
		return errors.New("email must not be empty")
	}
	if data.Password == "" {
		return errors.New("password must not be empty")
	}
	existingUser, err := s.userRepo.FindUserByEmail(data.Email)
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
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err
	}
	// send the email verification mail
	if s.cfg.Env != "test" {
		// generate the email verification token
		emailVerificationToken, err := s.jwtSvc.GenerateEmailVerificationToken(map[string]interface{}{
			"email": newUser.Email,
			"id":    newUser.ID.Hex(),
		})
		if err != nil {
			logrus.WithError(err).Error(err.Error())
			return fmt.Errorf("error generating email verification token: %s", err.Error())
		}
		emailVerificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.cfg.ClientURL, emailVerificationToken)

		if s.cfg.WithWorkers {
			sendVerificationMailTask, err := tasks.NewUserVerificationMailTask(
				newUser.Email,
				fmt.Sprintf("%s %s", newUser.Firstname, newUser.Lastname),
				"Verify your Kipa E-mail Address",
				struct {
					Name string
					URL  string
				}{
					Name: fmt.Sprintf("%s %s", newUser.Firstname, newUser.Lastname),
					URL:  emailVerificationLink,
				},
			)

			if err != nil {
				logrus.WithError(err).Error(err.Error())
			}

			s.queue.Add(sendVerificationMailTask, asynq.Queue("critical"))
			return nil
		} else {
			mailSvc := mailer.NewMailer(s.cfg)

			err = mailSvc.SendEmailVerificationMail(
				newUser.Email,
				fmt.Sprintf("%s %s", newUser.Firstname, newUser.Lastname),
				"Verify your Kipa E-mail Address",
				struct {
					URL string
				}{
					URL: emailVerificationLink,
				},
			)

			if err != nil {
				logrus.WithError(err).Error(err.Error())
			}
			logrus.Info("Successfully sent email verification link")
			return nil
		}

	}
	return nil
}

// Verify a user's email
// Accepts the email verification token
func (s *UserService) VerifyEmail(data dto.VerifyEmailInputDTO) error {
	claims, err := s.jwtSvc.DecodeToken(data.Token, s.cfg.EmailVerificationTokenSecretKey)
	if err != nil {
		return nil
	}
	// extract the user ID from the decoded JWT payload
	userID := claims.Payload["id"]

	// update the emailVerified status for the user
	user := &models.User{
		EmailVerified: true,
	}
	ID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return models.ErrInvalidObjectID
	}
	user.ID = ID
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// update the user's email verification status
	if err := s.userRepo.UpdateUser(user); err != nil {
		logrus.WithError(err).Error(err.Error())
		return fmt.Errorf("error verifying email: %s", err.Error())
	}

	return nil
}

// Returns the user data for the passed ID
// Accepts the user ID
func (s *UserService) FindUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.FindUserById(id)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

// Returns all the users in the database
func (s *UserService) FindAllUsers() ([]models.User, error) {
	users, err := s.userRepo.FindAllUsers()
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

// Update a user's details
// Accepts the user ID and the update user data
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
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

// Update a user's password
// Accepts the user ID and the new password
func (s *UserService) UpdateUserPassword(id string, data dto.UpdateUserPasswordInputDTO) error {
	// validation
	if utils.IsStringEmpty(id) {
		return ErrUserIDIsEmpty
	}
	if utils.IsStringEmpty(data.Password) {
		return ErrPasswordIsEmpty
	}
	// hash the new password
	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		Password: hashedPassword,
	}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ErrInvalidObjectID
	}
	user.ID = ID
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

// Delete a user from the database
// Accepts the user ID
func (s *UserService) DeleteUser(id string) error {
	if err := s.userRepo.DeleteUser(id); err != nil {
		return err
	}
	return nil
}
