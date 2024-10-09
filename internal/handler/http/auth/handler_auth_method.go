package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	// Parse login credentials from request body
	var loginDTO user.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	// Call the login service to authenticate user
	result, err := h.authService.Login(c, &loginDTO)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Login failed"})
	}

	// If login is successful, generate a cookie
	utils.GenerateCookie(c, result) // `result` can be a token or user ID depending on the logic

	// Return a success response
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "message": "Login successful"})
}
