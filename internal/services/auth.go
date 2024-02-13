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

type AuthService struct {
	userRepo repository.IUserRepository
	jwtSvc   jwt.IJwtService
	cfg      *config.Config
	queue    *queue.RedisQueue
}

type IAuthService interface {
	Login(data dto.LoginUserInputDTO) (*dto.LoginUserOutputDTO, error)
	RefreshToken(user *models.User) (*dto.RefreshTokenOutputDTO, error)
	ForgotPassword(data dto.ForgotPasswordInputDTO) error
	ResetPassword(data dto.ResetPasswordInputDTO) error
}

func NewAuthService(cfg *config.Config, userRepo repository.IUserRepository) IAuthService {
	jwtSvc := jwt.NewJwtService(cfg, userRepo)
	queue := queue.NewRedisQueue(cfg)
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		jwtSvc:   jwtSvc,
		queue:    queue,
	}
}

// error constants
var (
	ErrEmailIsEmpty    = errors.New("email cannot be empty")
	ErrPasswordIsEmpty = errors.New("password cannot be empty")
)

// Login user
// returns access and refresh token
func (s *AuthService) Login(data dto.LoginUserInputDTO) (*dto.LoginUserOutputDTO, error) {
	if utils.IsStringEmpty(data.Email) {
		return &dto.LoginUserOutputDTO{}, ErrEmailIsEmpty
	}
	if utils.IsStringEmpty(data.Password) {
		return &dto.LoginUserOutputDTO{}, ErrPasswordIsEmpty
	}
	// find the user assigned to input email
	user, err := s.userRepo.FindUserByEmail(data.Email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return &dto.LoginUserOutputDTO{}, err
	}
	if user == nil {
		return &dto.LoginUserOutputDTO{}, models.ErrUserNotFound
	}

	// compare passwords
	err = utils.ComparePasswordHash(data.Password, user.Password)
	if err != nil {
		return &dto.LoginUserOutputDTO{}, models.ErrIncorrectPassword
	}

	// if passwords match, generate access and refresh token
	payload := map[string]interface{}{"id": user.ID, "email": user.Email}
	accessToken, err := s.jwtSvc.GenerateAccessToken(payload)
	if err != nil {
		logrus.WithError(err).Error("error generating access token")
		return &dto.LoginUserOutputDTO{}, errors.New("error generating access token")
	}
	refreshToken, err := s.jwtSvc.GenerateRefreshToken(payload)
	if err != nil {
		logrus.WithError(err).Error("error generating refresh token")
		return &dto.LoginUserOutputDTO{}, errors.New("error generating refresh token")
	}
	return &dto.LoginUserOutputDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh token
// returns a refreshed access token
func (s *AuthService) RefreshToken(user *models.User) (*dto.RefreshTokenOutputDTO, error) {
	if utils.IsStringEmpty(user.Email) {
		return &dto.RefreshTokenOutputDTO{}, ErrEmailIsEmpty
	}
	if utils.IsStringEmpty(user.Password) {
		return &dto.RefreshTokenOutputDTO{}, ErrPasswordIsEmpty
	}
	payload := map[string]interface{}{"id": user.ID, "email": user.Email}
	accessToken, err := s.jwtSvc.GenerateAccessToken(payload)
	if err != nil {
		logrus.WithError(err).Error("error generating access token")
		return &dto.RefreshTokenOutputDTO{}, errors.New("error generating access token")
	}
	return &dto.RefreshTokenOutputDTO{
		AccessToken: accessToken,
	}, nil
}

// Forgot password
// Accepts the user's email
func (s *AuthService) ForgotPassword(data dto.ForgotPasswordInputDTO) error {
	// find the user
	user, err := s.userRepo.FindUserByEmail(data.Email)
	if err != nil {
		return err
	}
	// generate a reset password token
	payload := map[string]interface{}{
		"id":    user.ID.Hex(),
		"email": user.Email,
	}
	resetPasswordToken, err := s.jwtSvc.GenerateResetPasswordToken(payload)
	if err != nil {
		return fmt.Errorf("error generating reset password token: %s", err.Error())
	}
	// construct the reset password link
	resetPasswordLink := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.ClientURL, resetPasswordToken)
	// send the reset password link to the user's mail
	if s.cfg.Env != "test" {
		if s.cfg.WithWorkers {
			logrus.Info("sending with workers")
			sendResetPasswordMailTask, err := tasks.NewUserResetPasswordMailTask(
				user.Email,
				fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
				"Reset your Kipa Account Password",
				struct {
					Name string
					URL  string
				}{
					Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
					URL:  resetPasswordLink,
				},
			)

			if err != nil {
				logrus.WithError(err).Error(err.Error())
			}

			s.queue.Add(sendResetPasswordMailTask, asynq.Queue("critical"))
		} else {
			logrus.Info("sending without workers")
			mailSvc := mailer.NewMailer(s.cfg)

			err = mailSvc.SendResetPasswordMail(
				user.Email,
				fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
				"Reset your Kipa Account Password",
				struct {
					Name string
					URL  string
				}{
					Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
					URL:  resetPasswordLink,
				},
			)

			if err != nil {
				logrus.WithError(err).Error(err.Error())
				return err
			}
			logrus.Info("Successfully sent reset password link")
			return nil
		}
	}
	return nil
}

// Reset password
// Accepts the reset password token and the new password
func (s *AuthService) ResetPassword(data dto.ResetPasswordInputDTO) error {
	claims, err := s.jwtSvc.DecodeToken(data.Token, s.cfg.ResetPasswordTokenSecretKey)
	if err != nil {
		return nil
	}
	// extract the user ID from the decoded JWT payload
	userID := claims.Payload["id"]

	// hash the new password
	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return err
	}
	user := &models.User{
		Password: hashedPassword,
	}
	ID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return models.ErrInvalidObjectID
	}
	user.ID = ID
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// update the password
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}
