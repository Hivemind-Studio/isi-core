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
		return "", fmt.Errorf("missing salt environment variables")
	}

	saltedPassword := salt + password

	hash := sha256.New()
	hash.Write([]byte(saltedPassword))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ComparePassword(providedPassword, storedHash string) (bool, error) {
	hashedPassword, err := HashPassword(storedHash)

	if err != nil {
		return false, err
	}

	return hashedPassword == providedPassword, nil
}
