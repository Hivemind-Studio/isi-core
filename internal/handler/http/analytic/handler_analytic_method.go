package analytic

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetTotalRegistrant(c *fiber.Ctx) error {
	module := "Analytic Handler"
	functionName := "GetTotalRegistrant"

	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", "")

	res, err := h.getTotalRegistrant.Execute(c.Context())

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Analytic retrieved successfully",
		Data:    res,
	})
}
