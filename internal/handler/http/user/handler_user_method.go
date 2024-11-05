package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Create(c *fiber.Ctx) error {
	var newUser user.RegistrationDTO
	if err := c.BodyParser(&newUser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	err := newUser.ValidatePassword()

	if err != nil {
		return err
	}

	if err := validatorhelper.ValidateStruct(newUser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.userService.Create(c, &newUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "User created successfully",
			Data:    result,
		})
}
