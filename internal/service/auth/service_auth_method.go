package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Login(ctx *fiber.Ctx, body *user.LoginDTO) (result string, err error) {
	_, err = s.repoAuth.Login(ctx, body)

	if err != nil {
		return result, err
	}

	return s.repoAuth.Login(ctx, body)
}
