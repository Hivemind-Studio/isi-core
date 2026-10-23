package hash

import (
	"crypto/sha256"
	"encoding/hex"
	_ "errors"
	"fmt"
	"os"
)

func HashPassword(password string) (string, error) {
	salt := os.Getenv("PASSWORD_SALT")
	if salt == "" {
		return "", fmt.Errorf("salt is not set in environment variables")
	}

	saltedPassword := salt + password

	hash := sha256.New()
	hash.Write([]byte(saltedPassword))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ComparePassword(providedPassword, storedHash string) (bool, error) {
	hashedPassword, err := HashPassword(providedPassword)
	if err != nil {
		return false, err
	}

	return storedHash == hashedPassword, nil
}
