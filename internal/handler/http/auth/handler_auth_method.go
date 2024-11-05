package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	userRepo "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/cookie"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var loginDTO user.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	result, err := h.authService.Login(c, &loginDTO)
	if err != nil {
		return err
	}

	cookie.GenerateCookie(c, result)

	responseData := userRepo.LoginResponse{
		Email: result.Email,
		Name:  result.Name,
		Photo: result.Photo,
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data:    responseData,
		})
}
