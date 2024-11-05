package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

type serviceUserInterface interface {
	Create(ctx *fiber.Ctx, body *user.RegistrationDTO) (result *user.RegisterResponse, err error)
}
