package auth

import "github.com/golang-jwt/jwt/v4"

type IJwtService interface {
	GenerateToken(payload map[string]interface{}) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

// JWT custom claims
type JwtCustomClaims struct {
	payload map[string]interface{}
	jwt.StandardClaims
}
