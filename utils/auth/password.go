package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can not hash password %w", err)
	}
	return string(hash), nil
}

func VerifiedPassword(hash string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(candidatePassword))
}
