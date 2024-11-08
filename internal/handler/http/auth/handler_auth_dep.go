package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

type serviceAuthInterface interface {
	Login(ctx *fiber.Ctx, body *auth.LoginDTO) (result dto.UserDTO, err error)
}
