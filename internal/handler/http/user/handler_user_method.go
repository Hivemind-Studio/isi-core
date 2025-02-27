package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
)

func (h *Handler) Create(c *fiber.Ctx) error {
	var newUser auth.RegistrationStaffDTO
	if err := c.BodyParser(&newUser); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	if err := validatorhelper.ValidateStruct(newUser); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	err := h.createUserStaffUseCase.Execute(c.Context(), newUser)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "User created successfully",
		})
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	name := c.Query("name")
	email := c.Query("email")
	phoneNumber := c.Query("phone_number")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	pageParam := c.Query("page")
	perPageParam := c.Query("per_page")
	campaignId := c.Query("campaign_id")

	var start, end *time.Time
	if startDate != "" {
		parsedStart, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid end date format")
		}
		start = &parsedStart
	}
	if endDate != "" {
		parsedEnd, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid end date format")
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

	users, paginate, err := h.getUsersUseCase.Execute(c.Context(), name, email, phoneNumber, status, start, end,
		campaignId, page, perPage)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:     fiber.StatusOK,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: paginate,
	})
}

func (h *Handler) GetUserById(c *fiber.Ctx) error {
	paramId := c.Params("id")

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid user id")
	}

	res, err := h.getUserByIDUseCase.Execute(c.Context(), id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) UpdateStatusUser(c *fiber.Ctx) error {
	var payload user.SuspendDTO
	if err := c.BodyParser(&payload); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	err := h.updateUserStatusUseCase.Execute(c.Context(), payload.Ids, payload.UpdatedStatus)

	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Failed to update status users")
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Suspend users successfully",
		})
}

func (h *Handler) UpdateUserRole(c *fiber.Ctx) error {
	var payload user.UserRole
	paramId := c.Params("id")

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid user id")
	}

	if err := c.BodyParser(&payload); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}

	err = h.updateUserRoleCase.Execute(c.Context(), id, constant.GetRoleID(payload.Role))

	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Failed to update role users")
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Change role users successfully",
		})
}

func (h *Handler) SendChangeEmailVerification(c *fiber.Ctx) error {
	userSession, err := session.GetUserSession(c)
	if err != nil || userSession == nil {
		return fiber.ErrUnauthorized
	}

	go func() {
		err := h.sendChangeEmailVerificationUseCase.Execute(c.Context(), userSession.Email)
		if err != nil {
			log.Printf("Failed to send change email verification: %v", err)
		}
	}()
	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Failed to update role users")
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Change email verification successfully",
		})
}

func (h *Handler) SendConfirmationNewEmail(c *fiber.Ctx) error {
	var payload user.ChangeEmailDTO
	if err := c.BodyParser(&payload); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}
	userInfo := c.Locals("user_info")

	if userInfo == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	u, ok := userInfo.(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user data"})
	}

	email := u["email"].(string)

	err := h.sendConfirmationChangeNewEmailUseCase.Execute(c.Context(), payload.Token, payload.NewEmail, email)

	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Failed to update email users")
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Change email send to new email successfully",
		})
}

func (h *Handler) ConfirmChangeNewEmail(c *fiber.Ctx) error {
	var payload user.ConfirmChangeEmailDTO
	if err := c.BodyParser(&payload); err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Invalid Input")
	}
	userInfo := c.Locals("user_info")

	if userInfo == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	u, ok := userInfo.(map[string]interface{})
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user data"})
	}

	email := u["email"].(string)

	err := h.updateUserEmailUseCase.Execute(c.Context(), payload.Token, email)

	if err != nil {
		return httperror.Wrap(fiber.StatusBadRequest, err, "Failed to update email users")
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Change email users successfully",
		})
}
