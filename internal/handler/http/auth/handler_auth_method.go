package user

import (
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var loginDTO authdto.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	result, err := h.authService.Login(c.Context(), &loginDTO)
	if err != nil {
		return err
	}

	token, err := middleware.GenerateToken(
		middleware.User{
			Name:  result.Name,
			Email: result.Email,
			Role:  *result.Role,
		})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "login successful",
			Data: authdto.LoginResponse{
				Name:  result.Name,
				Email: result.Email,
				Role:  *result.Role,
				Photo: utils.SafeDereferenceString(result.Photo),
				Token: token,
			},
		})
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var signUp authdto.SignUpDTO
	if err := c.BodyParser(&signUp); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	err := h.authService.SignUp(c.Context(), &signUp)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Registration successful! Please check your email to verify your account.",
		})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var newUser authdto.RegistrationDTO
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, "Auth Handler", "Create",
		"", string(c.Body()))

	if err := c.BodyParser(&newUser); err != nil {
		logger.Print("error", requestId, "Auth Handler", "Create",
			"Invalid input", string(c.Body()))
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
		logger.Print("error", requestId, "Auth Handler", "Create",
			err.Error(), newUser)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Registration successful!",
			Data:    result,
		})
}
