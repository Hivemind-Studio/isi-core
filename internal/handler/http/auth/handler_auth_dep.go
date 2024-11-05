package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	userRepo "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/gofiber/fiber/v2"
)

type serviceAuthInterface interface {
	Login(ctx *fiber.Ctx, body *user.LoginDTO) (result userRepo.Login, err error)
}
