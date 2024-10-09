package utils

import (
	"crypto/sha256"
	"encoding/hex"
	_ "errors"
	"fmt"
	"os"
)

// GenerateSalt generates a random 16-byte salt
func HashPassword(password string) (string, error) {
	// Get the salt from environment variable
	salt := os.Getenv("PASSWORD_SALT")
	if salt == "" {
		return "", fmt.Errorf("salt is not set in environment variables")
	}

	// Combine password and salt
	saltedPassword := salt + password

	// Create SHA-256 hash
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))

	// Return the hashed password as a hex string
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ComparePassword(providedPassword, storedHash string) (bool, error) {
	fmt.Println("Password", providedPassword)
	fmt.Println("Hash", storedHash)
	// Hash the provided password using the same salt
	hashedPassword, err := HashPassword(storedHash)

	fmt.Println("HashedPassword", hashedPassword == storedHash)
	if err != nil {
		return false, err
	}

	// Compare the two hashes
	return hashedPassword == providedPassword, nil
}
