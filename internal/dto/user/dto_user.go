package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

type LoginDTO struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type RegistrationDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=4"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	Email           string `json:"email" validate:"required"`
}

type RegisterResponse struct {
	Name  string `json:"Name"`
	Email string `json:"email"`
}

func (p *RegistrationDTO) ValidatePassword() error {
	if p.Password != p.ConfirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password and confirm password do not match")
	}
	return nil
}
