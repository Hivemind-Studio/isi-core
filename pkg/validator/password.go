package validator

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func ValidatePassword(p *auth.RegistrationDTO) error {
	if p.Password != p.ConfirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password do not match")
	}
	return nil
}
