package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Create(c *fiber.Ctx) error {
	module := "Campaign Handler"
	functionName := "CreateCampaign"

	var requestBody campaign.DTO
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", string(c.Body()))

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	err := h.createCampaignUseCase.Execute(c.Context(), requestBody)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			err.Error(), requestBody)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Campaign created successfully",
		})
}
