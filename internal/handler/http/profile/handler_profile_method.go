package profile

import (
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) GetProfile(c *fiber.Ctx) error {
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

func (h *Handler) UpdateProfilePassword(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "UpdateProfilePassword"

	var requestBody authdto.UpdatePassword
	requestId := c.Locals("request_id").(string)

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if err := validator.ValidatePassword(&requestBody.Password, &requestBody.ConfirmPassword); err != nil {
		return err
	}

	err = h.updateProfilePassword.Execute(c.Context(), claims.ID, requestBody.CurrentPassword, requestBody.Password)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "change password successfully",
	})
}
