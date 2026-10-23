package createrole

import "github.com/gofiber/fiber/v2"

func (s *UseCase) Execute(ctx *fiber.Ctx, name string) (result string, err error) {
	return s.repoRole.CreateRole(ctx, name)
}
