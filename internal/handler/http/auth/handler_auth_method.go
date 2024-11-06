package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/cookie"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var loginDTO auth.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	result, err := h.authService.Login(c, &loginDTO)
	if err != nil {
		return err
	}

	cookie.GenerateCookie(c, result)

	res := auth.LoginResponse{
		Name:  result.Name,
		Email: *result.Email,
		Photo: utils.SafeDereferenceString(result.Photo),
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "login successful",
			Data:    res,
		})
}
