package user

import (
	"crypto/rand"
	"encoding/base64"
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"time"
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

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	//if err := h.verifyRegistrationTokenUseCase.Execute(c.Context(),
	//	requestBody.Email, requestBody.Token); err != nil {
	//	return err
	//}

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

func (h *Handler) GoogleLogin(c *fiber.Ctx) error {
	url := h.googleLoginUseCase.Execute(c)

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *Handler) generateStateOauthCookie() (string, fiber.Cookie) {
	expiration := time.Now().Add(24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := fiber.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HTTPOnly: true,
		Secure:   true,
	}

	return state, cookie
}

//
//func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
//	oauthState := c.Cookies("oauthstate")
//
//	if c.Query("state") != oauthState {
//		log.Println("Invalid OAuth Google state")
//		return c.Redirect("/", fiber.StatusTemporaryRedirect)
//	}
//
//	data, err := h.getUserDataFromGoogle(c.Query("code"))
//	if err != nil {
//		log.Println(err.Error())
//		return c.Redirect("/", fiber.StatusTemporaryRedirect)
//	}
//
//	// Return user data as a response
//	return c.JSON(fiber.Map{
//		"message":  "User data retrieved successfully",
//		"userData": string(data),
//	})
//}

//
//func (h *Handler) getUserDataFromGoogle(code string) ([]byte, error) {
//	token, err := googleOauthConfig.Exchange(context.Background(), code)
//	if err != nil {
//		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
//	}
//
//	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
//	if err != nil {
//		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
//	}
//	defer response.Body.Close()
//
//	contents, err := io.ReadAll(response.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read response: %s", err.Error())
//	}
//
//	return contents, nil
//}
