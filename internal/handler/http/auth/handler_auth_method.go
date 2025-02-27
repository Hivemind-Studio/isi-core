package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var loginDTO authdto.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	user, err := h.loginUseCase.Execute(c.Context(), &loginDTO)
	if err != nil {
		return err
	}

	token, err := h.setSession(c.Context(), user.ID, user.Email, user.Name, user.Role, user.Photo)
	if err != nil {
		return err
	}

	// Create a new cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour) // Set the cookie to expire in 24 hours
	cookie.HTTPOnly = true                          // Make the cookie accessible only via HTTP
	cookie.Secure = true                            // Set the cookie to be sent only over HTTPS
	cookie.Path = "/"                               // Set the cookie path to the root of the domain

	// Set the cookie in the response
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "login successful",
			Data: authdto.LoginResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  *user.Role,
				Photo: utils.SafeDereferenceString(user.Photo),
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

	if err := h.sendRegistrationEmailVerificationUseCase.Execute(c.Context(), requestBody.Email); err != nil {
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
	logger.Print(loglevel.INFO, requestId, module, functionName,
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

func (h *Handler) GoogleLogin(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "GoogleLogin"
	requestId := c.Locals("request_id").(string)
	logger.Print(loglevel.INFO, requestId, module, functionName,
		"", string(c.Body()))

	url := h.googleLoginUseCase.Execute(c)

	return c.Redirect(url)
}

func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "GoogleCallback"
	requestId := c.Locals("request_id").(string)
	logger.Print(loglevel.INFO, requestId, module, functionName,
		"", string(c.Body()))

	returnedState := c.Query("state")

	stateCookie := c.Cookies("oauth_state")
	if returnedState == "" || stateCookie == "" || returnedState != stateCookie {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	userData, err := h.googleCallbackUseCase.GetUserDataFromGoogle(c.Context(), c.Query("code"))
	if err != nil {
		return err
	}
	token, err := h.setSession(c.Context(), userData.ID, userData.Email,
		userData.Name, userData.Role, userData.Photo)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Login successful!",
			Data:    authdto.GoogleCallbackResponse{Token: token},
		})
}

func (h *Handler) setSession(c context.Context, id int64, email string, name string,
	role *string, photo *string) (token string, err error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	token = newUUID.String()
	userSession := session.Session{
		ID:    id,
		Email: email,
		Name:  name,
		Role:  utils.SafeDereferenceString(role),
		Photo: utils.SafeDereferenceString(photo),
	}
	err = h.sessionManager.CreateSession(c, "SESSION::"+token, userSession, time.Hour*1)
	if err != nil {
		return token, err
	}
	return token, nil
}
