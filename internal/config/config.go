package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// init is invoked on load
func init() {
	// loads values from .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load environment variables.")
	}
}

// config struct is used to store the loaded environment variables

type Config struct {
	Port                     string
	Env                      string
	AccessTokenJwtSecretKey  string
	RefreshTokenJwtSecretKey string
	AccessTokenJwtExpiresIn  string
	RefreshTokenJwtExpiresIn string
}

// New() creates a new Config struct with the loaded environment variables
func New() *Config {
	return &Config{
		Port:                     getEnv("PORT", "1323"),
		Env:                      getEnv("ENV", "development"),
		AccessTokenJwtSecretKey:  getEnv("ACCESS_TOKEN_JWT_SECRET_KEY", ""),
		RefreshTokenJwtSecretKey: getEnv("REFRESH_TOKEN_JWT_SECRET_KEY", ""),
		AccessTokenJwtExpiresIn:  getEnv("ACCESS_TOKEN_JWT_EXPIRES_IN", "15m"),
		RefreshTokenJwtExpiresIn: getEnv("REFRESH_TOKEN_JWT_EXPIRES_IN", "7d"),
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
	if value, err := strconv.Atoi(valueStr); err != nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool() reads the value for an environment variable 'key' as a boolean
// or returns a default value 'defaultValue'
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	// convert string value to a boolean
	if value, err := strconv.ParseBool(valueStr); err != nil {
		return value
	}
	return defaultValue
}

// getEnvAsFloat reads the value for an environment variable 'key' as a float
// or returns a default value 'defaultValue'
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	// convert string value to float32
	if value, err := strconv.ParseFloat(valueStr, 64); err != nil {
		return value
	}
	return defaultValue
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