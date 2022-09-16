package server

import (
	"errors"
	"fmt"
	"keeper/internal/auth"
	"keeper/internal/auth/auth_realm"
	"keeper/internal/config"
	"keeper/internal/models"
	"keeper/internal/repository"
	"keeper/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Middleware struct {
	Cfg              *config.Config
	UserRepository   repository.IUserRepository
	ApiKeyRepository repository.IAPIKeyRepository
	BucketRepository repository.IBucketRepository
}

var (
	apiKeyPermissionsCtxKey = "apikey_permissions"
	credTypeCtxKey          = "cred_type"
)

func NewMiddleware(cfg *config.Config, dbClient *mongo.Client) *Middleware {
	userRepository := repository.NewUserRepository(cfg, dbClient)
	apiKeyRepository := repository.NewAPIKeyRepository(cfg, dbClient)
	bucketRepository := repository.NewBucketRepository(cfg, dbClient)
	return &Middleware{
		Cfg:              cfg,
		UserRepository:   userRepository,
		ApiKeyRepository: apiKeyRepository,
		BucketRepository: bucketRepository,
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
		c.Set("user", authResponse.User)
		c.Set(apiKeyPermissionsCtxKey, authResponse.Permissions)
		c.Set(credTypeCtxKey, cred.Type)
		return next(c)
	}
}

// middleware for protecting refresh token route
func (m *Middleware) RequireRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		refreshToken := c.Request().Header.Get("x-refresh-token")
		// validate the auth header structure
		if len(strings.Split(refreshToken, ".")) != 3 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid refresh token structure",
			})
		}
		cred := &auth.Credential{
			Type: auth.CredentialTypeRefreshJWT,
			JWT:  refreshToken,
		}
		authRealm := auth_realm.NewAuthRealm(m.Cfg, m.ApiKeyRepository, m.UserRepository)
		authResponse, err := authRealm.Authenticate(cred)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": err,
			})
		}
		c.Set("user", authResponse.User)
		return next(c)
	}
}

// Middleware for protecting api key write bucket permission
func (m *Middleware) RequireAPIKeyBucketWritePermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to write bucket
			check := apiKeyPermissions.Contains(models.APIKeyPermissionWriteBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to write to bucket, %s permission is required.",
						models.APIKeyPermissionWriteBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key read bucket permission
func (m *Middleware) RequireAPIKeyBucketReadPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to read bucket
			check := apiKeyPermissions.Contains(models.APIKeyPermissionReadBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to read bucket, %s permission is required.",
						models.APIKeyPermissionReadBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key delete bucket permission
func (m *Middleware) RequireAPIKeyBucketDeletePermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to delete bucket
			check := apiKeyPermissions.Contains(models.APIKeyPermissionDeleteBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to delete bucket, %s permission is required.",
						models.APIKeyPermissionReadBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key read item permission
func (m *Middleware) RequireAPIKeyReadItemPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to read item
			check := apiKeyPermissions.Contains(models.APIKeyPermissionReadItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to read bucket item, %s permission is required.",
						models.APIKeyPermissionReadItem.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key write item permission
func (m *Middleware) RequireAPIKeyWriteItemPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to write item
			check := apiKeyPermissions.Contains(models.APIKeyPermissionWriteItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to write to bucket item, %s permission is required.",
						models.APIKeyPermissionWriteItem.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key delete item permission
func (m *Middleware) RequireAPIKeyDeleteItemPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to delete item
			check := apiKeyPermissions.Contains(models.APIKeyPermissionDeleteItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to delete bucket item, %s permission is required.",
						models.APIKeyPermissionDeleteItem.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key read user permission
func (m *Middleware) RequireAPIKeyReadUserPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to read user
			check := apiKeyPermissions.Contains(models.APIKeyPermissionReadUser)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to read user, %s permission is required.",
						models.APIKeyPermissionReadUser.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key write user permission
func (m *Middleware) RequireAPIKeyWriteUserPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to read user
			check := apiKeyPermissions.Contains(models.APIKeyPermissionWriteUser)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to write user, %s permission is required.",
						models.APIKeyPermissionWriteUser.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting api key write user permission
func (m *Middleware) RequireAPIKeyDeleteUserPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		if credType == auth.CredentialTypeAPIKey {
			apiKeyPermissions := c.Get(apiKeyPermissionsCtxKey).(models.APIKeyPermissionsList)
			// check if the api key permissions contains the permission to read user
			check := apiKeyPermissions.Contains(models.APIKeyPermissionDeleteUser)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to delete user, %s permission is required.",
						models.APIKeyPermissionDeleteUser.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket write access
func (m *Middleware) RequireBucketWriteAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicWriteBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to write to bucket, %s permission is required.",
						models.BucketPermissionPublicWriteBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket read access
func (m *Middleware) RequireBucketReadAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicReadBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to read bucket, %s permission is required.",
						models.BucketPermissionPublicReadBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket delete access
func (m *Middleware) RequireBucketDeleteAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicDeleteBucket)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to delete bucket, %s permission is required.",
						models.BucketPermissionPublicDeleteBucket.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket item write access
func (m *Middleware) RequireBucketItemWriteAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicWriteItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to write bucket item, %s permission is required.",
						models.BucketPermissionPublicWriteItem.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket item read access
func (m *Middleware) RequireBucketItemReadAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicReadItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to read bucket item, %s permission is required.",
						models.BucketPermissionPublicReadItem.String(),
					),
				})
			}
		}
		return next(c)
	}
}

// Middleware for protecting bucket item delete access
func (m *Middleware) RequireBucketItemDeleteAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bucketUID := c.Param("bucketUID")
		credType := c.Get(credTypeCtxKey).(auth.CredentialType)
		// check if the auth header credential type is an API Key
		if credType == auth.CredentialTypeAPIKey {
			// fetch the bucket from the uid
			bucket, err := m.BucketRepository.FindBucketByUID(bucketUID)
			if err != nil {
				logrus.WithError(err).Error("bucket does not exist")
				return c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Status: false,
					Error:  "bucket does not exist",
				})
			}
			check := bucket.Permissions.Contains(models.BucketPermissionPublicDeleteItem)
			if !check {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{
					Status: false,
					Error: fmt.Sprintf(
						"unable to delete bucket item, %s permission is required.",
						models.BucketPermissionPublicDeleteItem.String(),
					),
				})
			}
		}
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
		return &auth.Credential{
			Type: auth.CredentialTypeRefreshJWT,
			JWT:  refreshToken,
		}, nil
	default:
		return nil, fmt.Errorf("unknown credential authorization type: %s", credType)
	}

}
