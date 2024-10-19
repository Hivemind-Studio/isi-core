package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

func (s *AuthService) Login(ctx *fiber.Ctx, body *user.LoginDTO) (result string, err error) {
	return s.repoAuth.Login(ctx, body)
}
