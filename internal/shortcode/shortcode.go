package shortcode

import (
	"crypto/rand"
	"errors"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const DefaultLength = 8

func Generate() (string, error) {
	return GenerateShortCode(DefaultLength)
}

var ErrInvalidLength = errors.New("length cannot be negative or zero")

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}

	return string(b), nil
}
