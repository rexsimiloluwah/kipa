package jwt

import (
	"errors"
	"fmt"
	"keeper/internal/auth"
	"keeper/internal/config"
	"keeper/internal/repository"
	_ "keeper/pkg/log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type IJwtService interface {
	GenerateToken(payload map[string]interface{}, expiresIn string, secret string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	DecodeToken(tokenString string) (*JwtCustomClaims, error)
	GenerateAccessToken(payload map[string]interface{}) (string, error)
	GenerateRefreshToken(payload map[string]interface{}) (string, error)
	Authenticate(credential *auth.Credential) (*auth.AuthResponse, error)
}

// JWT custom claims
type JwtCustomClaims struct {
	Payload map[string]interface{}
	jwt.StandardClaims
}

type JwtService struct {
	Cfg            *config.Config
	Issuer         string
	UserRepository repository.IUserRepository
}

func NewJwtService(cfg *config.Config, userRepository repository.IUserRepository) IJwtService {
	return &JwtService{
		Cfg:            cfg,
		Issuer:         "rexsimiloluwa@gmail.com",
		UserRepository: userRepository,
	}
}

// Parses the expires in string (i.e. 7d, 24h, 60m etc.) to a Unix time format
func parseExpiresInTime(expiresInStr string) (int64, error) {
	durationStr := expiresInStr[:len(expiresInStr)-1]
	durationType := expiresInStr[len(expiresInStr)-1]
	// convert the expires in duration to a number
	durationNum, err := strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}
	var result int64
	switch durationType {
	case 'h':
		result = time.Now().Add(time.Hour * time.Duration(durationNum)).Unix()
	case 'm':
		result = time.Now().Add(time.Minute * time.Duration(durationNum)).Unix()
	case 'd':
		result = time.Now().Add(24 * time.Hour * time.Duration(durationNum)).Unix()
	}
	return result, nil
}

func (jwtSrv *JwtService) GenerateToken(payload map[string]interface{}, expiresIn string, secret string) (string, error) {
	// parse the expires in string
	expiresInUnix, err := parseExpiresInTime(expiresIn)
	if err != nil {
		logrus.Error("expires in is invalid")
		return "", err
	}
	// initialize the JWT claims
	jwtClaims := &JwtCustomClaims{
		payload,
		jwt.StandardClaims{
			ExpiresAt: expiresInUnix,
			Issuer:    jwtSrv.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// generated the encoded token using the secret key
	encodedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		logrus.Error("encoded token could not be generated: %s", err.Error())
		return "", err
	}

	logrus.Info("Successfully generated token!")
	return encodedToken, nil
}

func (jwtSrv *JwtService) GenerateAccessToken(payload map[string]interface{}) (string, error) {
	return jwtSrv.GenerateToken(payload, jwtSrv.Cfg.AccessTokenJwtExpiresIn, jwtSrv.Cfg.JwtSecretKey)
}

func (jwtSrv *JwtService) GenerateRefreshToken(payload map[string]interface{}) (string, error) {
	return jwtSrv.GenerateToken(payload, jwtSrv.Cfg.RefreshTokenJwtExpiresIn, jwtSrv.Cfg.JwtSecretKey)
}

func (jwtSrv *JwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Signing method validation
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key if otherwise
		return []byte(jwtSrv.Cfg.JwtSecretKey), nil
	})
}

func (jwtSrv *JwtService) DecodeToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSrv.Cfg.JwtSecretKey), nil
	})

	if err != nil {
		logrus.Errorf("error parsing token: %s", err.Error())
		return &JwtCustomClaims{}, fmt.Errorf("error parsing token: %s", err.Error())
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return &JwtCustomClaims{}, fmt.Errorf("invalid or expired jwt")
	}
	return claims, nil
}

func (jwtSrv *JwtService) Authenticate(credential *auth.Credential) (*auth.AuthResponse, error) {
	if credential.Type != auth.CredentialTypeJWT && credential.Type != auth.CredentialTypeRefreshJWT {
		return nil, errors.New("credential must be of jwt type")
	}

	tokenString := credential.JWT

	claims, err := jwtSrv.DecodeToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID := claims.Payload["id"]
	user, err := jwtSrv.UserRepository.FindUserById(userID.(string))

	if err != nil {
		return nil, err
	}
	// return the response
	return &auth.AuthResponse{
		AuthMode:   auth.CredentialTypeJWT,
		Credential: *credential,
		User:       user,
	}, nil
}
