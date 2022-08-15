package server

import (
	"errors"
	"fmt"
	"keeper/internal/auth"
	"keeper/internal/auth/auth_realm"
	"keeper/internal/config"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middleware struct {
	Cfg              *config.Config
	UserRepository   repository.IUserRepository
	ApiKeyRepository repository.IAPIKeyRepository
}

func NewMiddleware(cfg *config.Config, dbClient *mongo.Client) *Middleware {
	userRepository := repository.NewUserRepository(cfg, dbClient)
	apiKeyRepository := repository.NewAPIKeyRepository(cfg, dbClient)
	return &Middleware{
		Cfg:              cfg,
		UserRepository:   userRepository,
		ApiKeyRepository: apiKeyRepository,
	}
}

// middleware for protecting auth routes
func (m *Middleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cred, err := getAuthFromRequest(c)
		if err != nil {
			logrus.WithError(err).Error("failed to get auth from request")
			return echo.ErrUnauthorized
		}
		authRealm := auth_realm.NewAuthRealm(m.Cfg, m.ApiKeyRepository, m.UserRepository)
		authResponse, err := authRealm.Authenticate(cred)
		if err != nil {
			logrus.WithError(err).Errorf("error authenticating user with type: %s", cred.Type.String())
			return echo.ErrUnauthorized
		}
		logrus.Debug(authResponse)
		c.Set("user", authResponse.User)
		return next(c)
	}
}

// Returns the auth credential type and details from the request
func getAuthFromRequest(c echo.Context) (*auth.Credential, error) {
	authHeader := c.Request().Header.Get("Authorization")
	// split the auth header
	authInfo := strings.Split(authHeader, " ")
	// validate the auth header structure
	if len(authInfo) != 2 {
		err := errors.New("invalid auth header structure")
		return nil, err
	}

	// get the credential type
	credType := strings.ToUpper(authInfo[0])
	switch credType {
	case "BEARER":
		// extract the auth token
		authToken := authInfo[1]
		// check if the auth token is empty
		if utils.IsStringEmpty(authToken) {
			return nil, errors.New("token cannot be empty")
		}
		// check the prefix of the token to confirm if it is a JWT or API Key
		prefix := fmt.Sprintf("%s%s", utils.APIKeyPrefix, utils.APIKeySeperator)
		if strings.HasPrefix(authToken, prefix) {
			// it's an API Key
			return &auth.Credential{
				Type:   auth.CredentialTypeAPIKey,
				APIKey: authToken,
			}, nil
		}
		// else, it's a JWT type
		return &auth.Credential{
			Type: auth.CredentialTypeJWT,
			JWT:  authToken,
		}, nil
	case "X-REFRESH-TOKEN":
		// extract the refresh token
		refreshToken := authInfo[1]
		// check if the auth token is empty
		if utils.IsStringEmpty(refreshToken) {
			return nil, errors.New("token cannot be empty")
		}
		fmt.Print(refreshToken)
		return &auth.Credential{
			Type: auth.CredentialTypeRefreshJWT,
			JWT:  refreshToken,
		}, nil
	default:
		return nil, fmt.Errorf("unknown credential authorization type: %s", credType)
	}

}
