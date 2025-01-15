package validator

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ValidatePassword(password *string, confirmPassword *string) error {
	if password == nil || confirmPassword == nil {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password cannot be nil")
	}

	trimmedPassword := strings.TrimSpace(*password)
	trimmedConfirmPassword := strings.TrimSpace(*confirmPassword)

	if trimmedPassword != trimmedConfirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password do not match")
	}

	if len(trimmedPassword) < 8 || len(trimmedPassword) > 12 {
		return httperror.New(fiber.StatusBadRequest, "password must be between 8 and 12 characters")
	}

	if len(trimmedConfirmPassword) < 8 || len(trimmedConfirmPassword) > 12 {
		return httperror.New(fiber.StatusBadRequest, "confirm password must be between 8 and 12 characters")
	}

	return nil
}
