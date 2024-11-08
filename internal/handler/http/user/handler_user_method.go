package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func (h *Handler) Create(c *fiber.Ctx) error {
	var newUser auth.RegistrationDTO
	if err := c.BodyParser(&newUser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	err := validator.ValidatePassword(&newUser)

	if err != nil {
		return err
	}

	if err := validatorhelper.ValidateStruct(newUser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.userService.Create(c.Context(), &newUser)
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

func (h *Handler) GetUsers(c *fiber.Ctx) error {
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

	users, err := h.userService.GetUsers(c.Context(), name, email, start, end, page, perPage)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func (h *Handler) GetUserById(c *fiber.Ctx) error {
	paramId := c.Params("id")

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid user id")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    user,
	})
}

func (h *Handler) SuspendUsers(c *fiber.Ctx) error {
	var input struct {
		IDs []int64 `json:"ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	err := h.userService.SuspendUsers(c.Context(), input.IDs)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Suspend users successfully",
		})
}
