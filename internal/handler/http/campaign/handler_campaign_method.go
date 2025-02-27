package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func (h *Handler) Create(c *fiber.Ctx) error {
	module := "Campaign Handler"
	functionName := "CreateCampaign"

	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", string(c.Body()))

	var newCampaign campaign.DTO
	if err := c.BodyParser(&newCampaign); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	if err := validatorhelper.ValidateStruct(newCampaign); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	res, err := h.createCampaignUseCase.Execute(c.Context(), newCampaign)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Campaign created successfully",
			Data:    res,
		})
}

func (h *Handler) Get(c *fiber.Ctx) error {
	name := c.Query("name")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	pageParam := c.Query("page")
	perPageParam := c.Query("per_page")

	var start, end *time.Time
	if startDate != "" {
		parsedStart, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid start date format")
		}
		start = &parsedStart
	}
	if endDate != "" {
		parsedEnd, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid end date format")
		}
		end = &parsedEnd
	}

	page := int64(1)
	perPage := int64(10)

	if pageParam != "" {
		parsedPage, err := strconv.ParseInt(pageParam, 10, 64)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if perPageParam != "" {
		parsedPerPage, err := strconv.ParseInt(perPageParam, 10, 64)
		if err == nil && parsedPerPage > 0 {
			perPage = parsedPerPage
		}
	}

	campaigns, paginate, err := h.getCampaignUseCase.Execute(c.Context(), name, status, start, end, page, perPage)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:     fiber.StatusOK,
		Message:    "Campaign retrieved successfully",
		Data:       campaigns,
		Pagination: paginate,
	})
}

func (h *Handler) UpdateStatusCampaign(c *fiber.Ctx) error {
	module := "Campaign Handler"
	functionName := "UpdateStatusCampaign"

	paramId := c.Params("id")

	id, err := strconv.ParseInt(paramId, 10, 64)

	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", string(c.Body()))

	// TODO: Update ke patch status
	err = h.deleteCampaignUseCase.Execute(c.Context(), id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Campaign deleted successfully",
	})
}
