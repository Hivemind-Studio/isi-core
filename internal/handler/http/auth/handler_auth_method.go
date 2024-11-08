package user

import (
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var loginDTO authdto.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	result, err := h.authService.Login(c, &loginDTO)
	if err != nil {
		return err
	}

	token, err := middleware.GenerateToken(
		middleware.User{
			Name:  result.Name,
			Email: result.Email,
			Role:  *result.Role,
		})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "login successful",
			Data: authdto.LoginResponse{
				Name:  result.Name,
				Email: result.Email,
				Role:  *result.Role,
				Photo: utils.SafeDereferenceString(result.Photo),
				Token: token,
			},
		})
}
