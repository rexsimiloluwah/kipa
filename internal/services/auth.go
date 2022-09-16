package services

import (
	"errors"
	"keeper/internal/auth/jwt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"

	"github.com/sirupsen/logrus"
)

type AuthService struct {
	UserRepository repository.IUserRepository
	JwtSvc         jwt.IJwtService
	Cfg            *config.Config
}

type IAuthService interface {
	Login(data dto.LoginUserInputDTO) (*dto.LoginUserOutputDTO, error)
	RefreshToken(user *models.User) (*dto.RefreshTokenOutputDTO, error)
}

func NewAuthService(cfg *config.Config, userRepo repository.IUserRepository) IAuthService {
	jwtSvc := jwt.NewJwtService(cfg, userRepo)
	return &AuthService{
		UserRepository: userRepo,
		Cfg:            cfg,
		JwtSvc:         jwtSvc,
	}
}

// error constants
var (
	ErrEmailIsEmpty    = errors.New("email cannot be empty")
	ErrPasswordIsEmpty = errors.New("password cannot be empty")
)

// Login user
// returns access and refresh token
func (a *AuthService) Login(data dto.LoginUserInputDTO) (*dto.LoginUserOutputDTO, error) {
	if utils.IsStringEmpty(data.Email) {
		return &dto.LoginUserOutputDTO{}, ErrEmailIsEmpty
	}
	if utils.IsStringEmpty(data.Password) {
		return &dto.LoginUserOutputDTO{}, ErrPasswordIsEmpty
	}
	// find the user assigned to input email
	user, err := a.UserRepository.FindUserByEmail(data.Email)
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
	accessToken, err := a.JwtSvc.GenerateAccessToken(payload)
	if err != nil {
		logrus.WithError(err).Error("error generating access token")
		return &dto.LoginUserOutputDTO{}, errors.New("error generating access token")
	}
	refreshToken, err := a.JwtSvc.GenerateRefreshToken(payload)
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
func (a *AuthService) RefreshToken(user *models.User) (*dto.RefreshTokenOutputDTO, error) {
	if utils.IsStringEmpty(user.Email) {
		return &dto.RefreshTokenOutputDTO{}, ErrEmailIsEmpty
	}
	if utils.IsStringEmpty(user.Password) {
		return &dto.RefreshTokenOutputDTO{}, ErrPasswordIsEmpty
	}
	payload := map[string]interface{}{"id": user.ID, "email": user.Email}
	accessToken, err := a.JwtSvc.GenerateAccessToken(payload)
	if err != nil {
		logrus.WithError(err).Error("error generating access token")
		return &dto.RefreshTokenOutputDTO{}, errors.New("error generating access token")
	}
	return &dto.RefreshTokenOutputDTO{
		AccessToken: accessToken,
	}, nil
}
