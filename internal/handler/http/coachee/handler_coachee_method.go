package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func (h *Handler) GetCoaches(c *fiber.Ctx) error {
	name := c.Query("name")
	email := c.Query("email")
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

	users, err := h.getCoacheesUseCase.Execute(c.Context(), name, email, start, end, page, perPage)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}
