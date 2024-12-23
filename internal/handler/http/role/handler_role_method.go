package role

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var dto struct {
		Name string `json:"name"`
	}

	if err := utils.ParseBody(c, &dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.WebResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Invalid request body",
			})
	}

	result, err := h.createRoleUseCase.Execute(c, dto.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.WebResponse{
				Status:  fiber.StatusBadRequest,
				Message: err.Error(),
			})
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    result,
	})
}
