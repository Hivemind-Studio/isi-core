package profile

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) GetProfileUser(c *fiber.Ctx) error {
	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	user, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	res, err := h.getProfileUser.Execute(c.Context(), user.ID)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) changePassword(c *fiber.Ctx) error {
	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	_, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "change password successfully",
	})
}
