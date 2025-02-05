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

	result, err := h.loginUseCase.Execute(c.Context(), &loginDTO)
	if err != nil {
		return err
	}

	token, err := middleware.GenerateToken(
		middleware.User{
			ID:    result.ID,
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
				ID:    result.ID,
				Name:  result.Name,
				Email: result.Email,
				Role:  *result.Role,
				Photo: utils.SafeDereferenceString(result.Photo),
				Token: token,
			},
		})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "CreateUser"

	var requestBody authdto.RegistrationDTO
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", string(c.Body()))

	c.Locals("user_info")

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	if err := h.verifyRegistrationTokenUseCase.Execute(c.Context(),
		requestBody.Email, requestBody.Token); err != nil {
		return err
	}

	if err := validator.ValidatePassword(&requestBody.Password, &requestBody.ConfirmPassword); err != nil {
		return err
	}

	if err := validatorhelper.ValidateStruct(requestBody); err != nil {
		return httperror.New(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.createUserUseCase.Execute(c.Context(), &requestBody)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			err.Error(), requestBody)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Registration successful!",
			Data:    result,
		})
}

func (h *Handler) EmailVerification(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "EmailVerification"

	var requestBody authdto.EmailVerificationDTO
	requestId := c.Locals("request_id").(string)

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	if err := h.sendEmailVerificationUseCase.Execute(c.Context(), requestBody.Email); err != nil {
		logger.Print("error", requestId, module, functionName,
			err.Error(), requestBody)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Email verification sent",
		})
}

func (h *Handler) VerifyEmailToken(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "VerifyEmailToken"

	var requestBody authdto.EmailVerificationDTO
	requestId := c.Locals("request_id").(string)

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	if err := h.sendEmailVerificationUseCase.Execute(c.Context(), requestBody.Email); err != nil {
		logger.Print("error", requestId, module, functionName,
			err.Error(), requestBody)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Email verification sent",
		})
}

func (h *Handler) PatchPassword(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "UpdatePassword"

	var requestBody authdto.UpdatePasswordRegistration
	requestId := c.Locals("request_id").(string)

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	if err := validator.ValidatePassword(&requestBody.Password, &requestBody.ConfirmPassword); err != nil {
		return err
	}

	err := h.updateCoachPasswordUseCase.Execute(c.Context(), requestBody.Password,
		requestBody.ConfirmPassword, requestBody.Token)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Update Password successful!",
		})
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "ForgotPassword"

	var requestBody authdto.ForgotPasswordDTO
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName,
		"", string(c.Body()))

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	err := h.forgotPasswordUseCase.Execute(c.Context(), requestBody.Email)

	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(
		response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "Forgot Password successful!",
		})
}
