package user

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	validatorhelper "github.com/Hivemind-Studio/isi-core/pkg/translator"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
	"time"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "Login"
	requestId := c.Locals("request_id").(string)

	var loginDTO authdto.LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	user, err := h.loginUseCase.Execute(c.Context(), &loginDTO)
	if err != nil {
		return err
	}

	token, err := h.generateAndSaveSessionToken(c.Context(), user.ID, user.Email, user.Name, user.Role, user.Photo)
	if err != nil {
		return err
	}

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("token created: %s", token), string(c.Body()))

	if user.RoleID == nil {
		logger.Print("error", requestId, module, functionName,
			fmt.Sprintf("user does not have role", user), string(c.Body()))
		return httperror.New(fiber.StatusUnauthorized, "User has no role")

	}

	err = h.setCookieByRole(c, user.RoleID, token)
	if err != nil {
		return err
	}

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("cookie set created: %s", token), string(c.Body()))

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

func (h *Handler) setCookieByRole(c *fiber.Ctx, roleId *int64, token string) error {
	origin := c.Get("Origin")

	if constant.IsDashboardUser(*roleId) {
		backoffice := "backoffice"
		if strings.Contains(origin, backoffice) {
			return httperror.New(fiber.StatusForbidden, "Invalid user role for dashboard application")
		}
		setCookie(c, token, backoffice)
	} else if constant.IsDashboardUser(*roleId) {
		dashboard := "dashboard"
		if strings.Contains(origin, dashboard) {
			return httperror.New(fiber.StatusForbidden, "Invalid user role for backoffice application")
		}
		setCookie(c, token, "backoffice")
	}
	return nil
}

func setCookie(c *fiber.Ctx, token string, domain string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour) // Set the cookie to expire in 24 hours
	cookie.HTTPOnly = true                          // Make the cookie accessible only via HTTP
	cookie.Secure = true                            // Set the cookie to be sent only over HTTPS
	cookie.Path = "/"
	cookie.SameSite = "None"
	cookie.Domain = fmt.Sprintf("%s+.inspirasisatu.com", domain)

	c.Cookie(cookie)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "CreateUser"
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

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

	requestBody.IPAddress = ipAddress
	requestBody.UserAgent = userAgent

	result, err := h.createUserUseCase.Execute(c.Context(), &requestBody)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			err.Error(), requestBody)
		return err
	}

	if requestBody.CampaignId != nil {
		userCampaign := campaign.UserCampaign{
			Email:      result.Email,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CampaignId: *requestBody.CampaignId,
		}
		err = h.createUserCampaign.Execute(c.Context(), userCampaign)

		if err != nil {
			return nil
		}
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

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	if err := h.sendRegistrationEmailVerificationUseCase.Execute(ctx, requestBody.Email); err != nil {
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

	stateCookie := c.Cookies("oauthstate")
	if returnedState == "" || stateCookie == "" || returnedState != stateCookie {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	userData, err := h.googleCallbackUseCase.GetUserDataFromGoogle(c.Context(), c.Query("code"))
	if err != nil {
		return err
	}
	token, err := h.generateAndSaveSessionToken(c.Context(), userData.ID, userData.Email,
		userData.Name, userData.Role, userData.Photo)
	if err != nil {
		return err
	}

	err = h.setCookieByRole(c, userData.RoleID, token)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Login successful!",
			Data: authdto.GoogleCallbackResponse{
				Token: token,
				ID:    userData.ID,
				Name:  userData.Name,
				Role:  userData.Role,
				Photo: userData.Photo,
			},
		})
}

func (h *Handler) generateAndSaveSessionToken(c context.Context, id int64, email string, name string,
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
