package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// init is invoked on load
// func init() {
// 	// loads values from .env
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("Could not load environment variables.")
// 	}
// }

// config struct is used to store the loaded environment variables

type Config struct {
	Port                            string
	Env                             string
	JwtSecretKey                    string
	EmailVerificationTokenSecretKey string
	ResetPasswordTokenSecretKey     string
	AccessTokenJwtExpiresIn         string
	RefreshTokenJwtExpiresIn        string
	EmailVerificationTokenExpiresIn string
	ResetPasswordTokenExpiresIn     string
	DbName                          string
	DbHost                          string
	DbPort                          string
	DbUser                          string
	DbPassword                      string
	MongoDbProdConnUri              string
	MongoDbTestConnUri              string
	SendgridAPIKey                  string
	ClientURL                       string
	RedisHost                       string
	RedisPort                       string
	RedisProdUri                    string
	WithWorkers                     bool
	EmailSenderName                 string
	EmailSenderAddr                 string
}

// New() creates a new Config struct with the loaded environment variables
func New() *Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Could not load environment variables.")
	}
	return &Config{
		Port:                            getEnv("PORT", "1323"),
		Env:                             getEnv("ENV", "development"),
		JwtSecretKey:                    getEnv("JWT_SECRET_KEY", ""),
		EmailVerificationTokenSecretKey: getEnv("EMAIL_VERIFICATION_TOKEN_SECRET_KEY", ""),
		ResetPasswordTokenSecretKey:     getEnv("RESET_PASSWORD_TOKEN_SECRET_KEY", ""),
		AccessTokenJwtExpiresIn:         getEnv("ACCESS_TOKEN_JWT_EXPIRES_IN", "15m"),
		RefreshTokenJwtExpiresIn:        getEnv("REFRESH_TOKEN_JWT_EXPIRES_IN", "7d"),
		EmailVerificationTokenExpiresIn: getEnv("EMAIL_VERIFICATION_TOKEN_EXPIRES_IN", "24h"),
		ResetPasswordTokenExpiresIn:     getEnv("RESET_PASSWORD_TOKEN_EXPIRES_IN", "7d"),
		DbHost:                          getEnv("MONGODB_HOST", ""),
		DbPort:                          getEnv("MONGODB_PORT", ""),
		DbUser:                          getEnv("MONGODB_USER", ""),
		DbPassword:                      getEnv("MONGODB_PASSWORD", ""),
		DbName:                          getEnv("MONGODB_NAME", "keeper"),
		MongoDbProdConnUri:              getEnv("MONGODB_CONN_URI", ""),
		MongoDbTestConnUri:              getEnv("MONGODB_TEST_CONN_URI", ""),
		SendgridAPIKey:                  getEnv("SENDGRID_API_KEY", ""),
		ClientURL:                       getEnv("CLIENT_URL", ""),
		RedisHost:                       getEnv("REDIS_HOST", ""),
		RedisPort:                       getEnv("REDIS_PORT", ""),
		RedisProdUri:                    getEnv("REDIS_PROD_URI", ""),
		WithWorkers:                     getEnvAsBool("WITH_WORKERS", true),
		EmailSenderName:                 getEnv("EMAIL_SENDER_NAME", "kipa"),
		EmailSenderAddr:                 getEnv("EMAIL_SENDER_ADDR", "rexsimiloluwa@gmail.com"),
	}
}

// NewTest() creates test environment variables
func NewTest() *Config {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Could not load environment variables.")
	// }
	return &Config{
		Port:                            getEnv("PORT", "1323"),
		Env:                             getEnv("ENV", "test"),
		JwtSecretKey:                    getEnv("JWT_SECRET_KEY", ""),
		EmailVerificationTokenSecretKey: getEnv("EMAIL_VERIFICATION_TOKEN_SECRET_KEY", ""),
		ResetPasswordTokenSecretKey:     getEnv("RESET_PASSWORD_TOKEN_SECRET_KEY", ""),
		AccessTokenJwtExpiresIn:         getEnv("ACCESS_TOKEN_JWT_EXPIRES_IN", "15m"),
		RefreshTokenJwtExpiresIn:        getEnv("REFRESH_TOKEN_JWT_EXPIRES_IN", "7d"),
		EmailVerificationTokenExpiresIn: getEnv("EMAIL_VERIFICATION_TOKEN_EXPIRES_IN", ""),
		ResetPasswordTokenExpiresIn:     getEnv("RESET_PASSWORD_TOKEN_EXPIRES_IN", ""),
		DbHost:                          getEnv("MONGODB_HOST", ""),
		DbName:                          getEnv("MONGODB_NAME", "keeper-go-test"),
		DbUser:                          getEnv("MONGODB_USER", ""),
		DbPort:                          getEnv("MONGODB_PORT", ""),
		DbPassword:                      getEnv("MONGODB_PASSWORD", ""),
		MongoDbTestConnUri:              getEnv("MONGODB_TEST_CONN_URI", ""),
		ClientURL:                       getEnv("CLIENT_URL", ""),
		RedisHost:                       getEnv("REDIS_HOST", ""),
		RedisPort:                       getEnv("REDIS_PORT", ""),
		RedisProdUri:                    getEnv("REDIS_PROD_URI", ""),
	}
}

// Helper function to read env variables
// getEnv() reads the value for an environment variable 'key'
// or returns a default value 'defaultValue'
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt() reads the value for an environment variable 'key' as an integer
// or returns a default value 'defaultValue'
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	// convert the string value to an integer
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool() reads the value for an environment variable 'key' as a boolean
// or returns a default value 'defaultValue'
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	// convert string value to a boolean
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsFloat reads the value for an environment variable 'key' as a float
// or returns a default value 'defaultValue'
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	// convert string value to float32
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsSlice reads the value for an environment variable 'key' as a slice
// or returns a default value 'defaultValue'
func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}
