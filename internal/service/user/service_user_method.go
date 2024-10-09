package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) GetTest(ctx *fiber.Ctx, id int) (result string, err error) {
	return s.repoUser.GetTest(ctx, id)
}

func (s *Service) Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterResponse, err error) {
	return s.repoUser.Create(ctx, body)
}
