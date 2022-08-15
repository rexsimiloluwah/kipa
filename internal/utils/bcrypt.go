package utils

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Hashing password using bcrypt
func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, 12)

	if err != nil {
		return "", errors.New("error hashing password")
	}

	return string(hashedPassword), nil
}

// Comparing hashed password
func ComparePasswordHash(password string, hash string) error {
	passwordBytes := []byte(password)
	hashedPasswordBytes := []byte(hash)

	if err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes); err != nil {
		return errors.New("error comparing password hash")
	}
	return nil
}
