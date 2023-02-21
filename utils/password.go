package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	// dari database, password input
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/*
	change password:

	- input password lama
	- input password baru
	- validasi password lama == password di database
	- kalo sam, update password baru
	- logout, frontend
*/
