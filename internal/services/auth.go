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
	"go.mongodb.org/mongo-driver/mongo"
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

func NewAuthService(cfg *config.Config, dbClient *mongo.Client) IAuthService {
	userRepository := repository.NewUserRepository(cfg, dbClient)
	jwtSvc := jwt.NewJwtService(cfg, userRepository)
	return &AuthService{
		UserRepository: userRepository,
		Cfg:            cfg,
		JwtSvc:         jwtSvc,
	}
}

// Login user
// returns access and refresh token
func (a *AuthService) Login(data dto.LoginUserInputDTO) (*dto.LoginUserOutputDTO, error) {
	// find the user assigned to input email
	user, err := a.UserRepository.FindUserByEmail(data.Email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return &dto.LoginUserOutputDTO{}, err
	}
	if user == nil {
		return &dto.LoginUserOutputDTO{}, errors.New("user not found")
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
