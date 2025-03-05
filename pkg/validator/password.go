package validator

import (
	"regexp"
	"strings"

	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func ValidatePassword(password *string, confirmPassword *string) error {
	if password == nil || confirmPassword == nil {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password cannot be nil")
	}

	trimmedPassword := strings.TrimSpace(*password)
	trimmedConfirmPassword := strings.TrimSpace(*confirmPassword)

	if len(trimmedPassword) < 8 || len(trimmedPassword) > 12 {
		return httperror.New(fiber.StatusBadRequest, "password must be between 8 and 12 characters")
	}

	uppercasePattern := regexp.MustCompile(`[A-Z]`)
	if !uppercasePattern.MatchString(trimmedPassword) {
		return httperror.New(fiber.StatusBadRequest, "password must contain at least one uppercase letter")
	}

	numberPattern := regexp.MustCompile(`[0-9]`)
	if !numberPattern.MatchString(trimmedPassword) {
		return httperror.New(fiber.StatusBadRequest, "password must contain at least one number")
	}

	specialPattern := regexp.MustCompile(`[\W_]`)
	if !specialPattern.MatchString(trimmedPassword) {
		return httperror.New(fiber.StatusBadRequest, "password must contain at least one special character")
	}

	if trimmedPassword != trimmedConfirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password do not match")
	}

	return nil
}
