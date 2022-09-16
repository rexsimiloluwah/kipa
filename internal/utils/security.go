package utils

import (
	"math/rand"
	"strings"

	"github.com/dchest/uniuri"
)

//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const (
	letterIdxBits   = 6                    // 6 bits to represent a letter index
	letterIdxMask   = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax    = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	APIKeyPrefix    = "KP"
	APIKeySeperator = "."
)

// Checks if a string is empty or not
func IsStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Generate random bytes of length n
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// // Generate random string of length n (using runes) - least performant
// func GenerateRandomStringUsingRunes(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }

// // Generate random string of length n (using string and rand.Int63()) - slightly more performant
// func GenerateRandomStringUsingRandInt63(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
// 	}
// 	return string(b)
// }

// Generate Random string of length n (main algo)
// source: - https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b), nil
}

// Generate a secret of length n
func GenerateSecret() (string, error) {
	return GenerateRandomString(24)
}

// Generates an API Key, returns the mask and generated API Key
func GenerateAPIKey() (string, string) {
	mask := uniuri.NewLen(16)
	key := uniuri.NewLen(48)

	var apiKey strings.Builder
	apiKey.WriteString(APIKeyPrefix)
	apiKey.WriteString(APIKeySeperator)
	apiKey.WriteString(mask)
	apiKey.WriteString(APIKeySeperator)
	apiKey.WriteString(key)

	return mask, apiKey.String()
}
