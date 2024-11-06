package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	userRepo "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/gofiber/fiber/v2"
)

type serviceAuthInterface interface {
	Login(ctx *fiber.Ctx, body *auth.LoginDTO) (result userRepo.Cookie, err error)
}
