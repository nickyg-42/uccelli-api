package helpers

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random 16-character string using uppercase letters and numbers.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	for i := range bytes {
		randomByte, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		bytes[i] = charset[randomByte.Int64()]
	}
	return string(bytes), nil
}
