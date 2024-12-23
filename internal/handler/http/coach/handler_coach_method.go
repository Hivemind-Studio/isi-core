package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func (h *Handler) GetCoaches(c *fiber.Ctx) error {
	name := c.Query("name")
	email := c.Query("useremail")
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

	users, err := h.getCoachUseCase.Execute(c.Context(), name, email, start, end, page, perPage)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func (h *Handler) CreateCoach(c *fiber.Ctx) error {
	var newCoach coach.CreateCoachDTO
	if err := c.BodyParser(&newCoach); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	if err := validatorhelper.ValidateStruct(newCoach); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	err := h.createCoachUseCase.Execute(c.Context(), newCoach)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Coach created successfully",
		})
}
