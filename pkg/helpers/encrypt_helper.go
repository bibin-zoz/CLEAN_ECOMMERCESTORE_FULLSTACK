package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func VerifyPassword(password, hashedPassword string) error {
	fmt.Println(password, hashedPassword)

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err // Passwords don't match
	}
	return nil // Passwords matchsss
}
