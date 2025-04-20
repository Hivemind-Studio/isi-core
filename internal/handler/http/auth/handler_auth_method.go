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
	"log"
	"net/http"
	"os"
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

	token, err := h.generateAndSaveSessionToken(c.Context(), user.ID, user.Email, user.Name, user.Role, user.Photo, requestId)
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
				Role:  utils.SafeDereferenceString(user.Role),
				Photo: utils.SafeDereferenceString(user.Photo),
				Token: token,
			},
		})
}

func (h *Handler) setCookieByRole(c *fiber.Ctx, roleId *int64, token string) error {
	module := "Auth Handler"
	functionName := "setCookieByRole"
	requestId := c.Locals("request_id").(string)

	if roleId == nil {
		logger.Print("error", requestId, module, functionName, "Role ID is missing", nil)
		return httperror.New(fiber.StatusBadRequest, "Role ID is required")
	}

	origin := c.Get("Origin")
	logger.Print("info", requestId, module, functionName, fmt.Sprintf("Request Origin: %s", origin), nil)

	if constant.PRODUCTION == h.env {
		logger.Print("info", requestId, module, functionName, "Production environment detected", nil)
		return h.setProductionCookie(c, *roleId, token, origin, requestId)
	}

	if constant.IsDashboardUser(*roleId) {
		setCookie(c, token, "", constant.TOKEN_ACCESS_DASHBOARD)
	} else if constant.IsBackofficeUser(*roleId) {
		setCookie(c, token, "", constant.TOKEN_ACCESS_BACKOFFICE)
	}
	logger.Print("info", requestId, module, functionName, "Cookie set for non-production environment", nil)
	return nil
}

func (h *Handler) setProductionCookie(c *fiber.Ctx, roleId int64, token, origin string, requestId string) error {
	module := "Auth Handler"
	functionName := "setProductionCookie"
	const domain = ".inspirasisatu.com"

	logger.Print("info", requestId, module, functionName, fmt.Sprintf("Processing role: %d, Origin: %s", roleId, origin), nil)

	if constant.IsDashboardUser(roleId) {
		if isOriginBackoffice(origin) {
			logger.Print("error", requestId, module, functionName, "Dashboard user detected but request came from backoffice", nil)
			return httperror.New(fiber.StatusForbidden, "Invalid user role for dashboard application")
		}
		setCookie(c, token, domain, constant.TOKEN_ACCESS_DASHBOARD)
		logger.Print("info", requestId, module, functionName, fmt.Sprintf("Cookie set for Dashboard "), nil)
	} else {
		if isOriginDashboard(origin) {
			logger.Print("error", requestId, module, functionName, "Backoffice user detected but request came from dashboard", nil)
			return httperror.New(fiber.StatusForbidden, "Invalid user role for backoffice application")
		}
		setCookie(c, token, domain, constant.TOKEN_ACCESS_BACKOFFICE)
		logger.Print("info", requestId, module, functionName, fmt.Sprintf("Cookie set for Backoffice "), nil)
	}
	return nil
}

func isOriginDashboard(origin string) bool {
	return strings.Contains(origin, "dashboard")
}

func isOriginBackoffice(origin string) bool {
	return strings.Contains(origin, "backoffice")
}

func setCookie(c *fiber.Ctx, token string, domain string, sessionCookieName string) {
	cookie := new(fiber.Cookie)
	cookie.Name = sessionCookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour) // Set the cookie to expire in 24 hours
	cookie.HTTPOnly = true                          // Make the cookie accessible only via HTTP
	cookie.Secure = true                            // Set the cookie to be sent only over HTTPS
	cookie.Path = "/"
	cookie.SameSite = "None"
	if domain != "" {
		cookie.Domain = domain
	}

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
		logger.Print(loglevel.ERROR, requestId, module, functionName,
			"Registration token is not valid", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Registration token is not valid")
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

	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"Incoming request", string(c.Body()))

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")
	campaignId := c.Query("campaign_id")
	returnedState := c.Query("state")
	stateCookie := c.Cookies("oauthstate")

	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"State Verification", fmt.Sprintf("ReturnedState: %s, StateCookie: %s", returnedState, stateCookie))

	if returnedState == "" || stateCookie == "" || returnedState != stateCookie {
		logger.Print(loglevel.WARN, requestId, module, functionName,
			"OAuth State Mismatch", "Invalid OAuth state detected")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	code := c.Query("code")
	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"Exchanging code", fmt.Sprintf("Code: %s", code))

	userData, err := h.googleCallbackUseCase.GetUserDataFromGoogle(c.Context(), code)
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, module, functionName,
			"GetUserDataFromGoogle failed", err.Error())

		return err
	}

	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"User data retrieved", fmt.Sprintf("ID: %d, Email: %s, Name: %s", userData.ID,
			userData.Email, userData.Name))

	token, err := h.generateAndSaveSessionToken(c.Context(), userData.ID, userData.Email, userData.Name,
		userData.Role, userData.Photo, "")
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, module, functionName,
			"Session token generation failed", err.Error())

		return err
	}

	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"Session token generated", fmt.Sprintf("Token: %s", token))

	err = h.setCookieByRole(c, userData.RoleID, token)
	if err != nil {
		logger.Print(loglevel.ERROR, requestId, module, functionName,
			"Setting cookie by role failed", err.Error())

		return err
	}

	logger.Print(loglevel.DEBUG, requestId, module, functionName,
		"Login successful", fmt.Sprintf("UserID: %d, Role: %v", userData.ID, userData.Role))

	if campaignId != "" {
		userCampaign := campaign.UserCampaign{
			Email:      userData.Email,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CampaignId: campaignId,
		}
		err = h.createUserCampaign.Execute(c.Context(), userCampaign)

		if err != nil {
			return nil
		}
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

func (h *Handler) generateAndSaveSessionToken(c context.Context, id int64, email string, name string, role *string, photo *string, requestId string) (token string, err error) {
	module := "Auth Handler"
	functionName := "generateAndSaveSessionToken"
	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf(
			"Generating session token with params: ID=%d, Email=%s, Name=%s, Role=%v, Photo=%v",
			id, email, name, role, photo,
		),
		nil,
	)

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

	logger.Print("info", requestId, module, functionName,
		"session created", userSession,
	)
	err = h.sessionManager.CreateSession(c, "SESSION::"+token, userSession, time.Hour*1)
	if err != nil {
		return token, httperror.New(http.StatusInternalServerError,
			fmt.Sprintf("Failed to create session token, %v", err))
	}
	return token, nil
}

func (h *Handler) DeleteSession(c *fiber.Ctx) error {
	module := "Auth Handler"
	functionName := "DeleteSession"
	requestId, _ := c.Locals("request_id").(string) // Avoid panic if nil

	origin := c.Get("Origin")
	environment := os.Getenv("ENVIRONMENT") // Ensure this is set correctly
	appOrigin := c.Get("X-App-Origin")      // Ensure this is passed in request headers
	var cookieName string

	if environment == constant.PRODUCTION {
		if utils.IsOriginDashboard(origin) || appOrigin == constant.DASHBOARD {
			cookieName = constant.TOKEN_ACCESS_DASHBOARD
		} else if utils.IsOriginBackoffice(origin) || appOrigin == constant.BACKOFFICE {
			cookieName = constant.TOKEN_ACCESS_BACKOFFICE
		} else {
			log.Println("Unauthorized access: Invalid token origin in production")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token origin"})
		}
	} else {
		if appOrigin == constant.DASHBOARD {
			cookieName = constant.TOKEN_ACCESS_DASHBOARD
		} else if appOrigin == constant.BACKOFFICE {
			cookieName = constant.TOKEN_ACCESS_BACKOFFICE
		} else {
			log.Println("Unauthorized access: Invalid token origin in non-production")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token origin"})
		}
	}

	// Get the token from the cookie
	token := c.Cookies(cookieName)
	if token == "" {
		log.Println("Unauthorized access: No session token found in cookies")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No session"})
	}

	logger.Print("info", requestId, module, functionName, "deleting session:", token)

	// Delete the session from Redis (or other session store)
	err := h.sessionManager.DeleteSession(c.Context(), "SESSION::"+token)
	if err != nil {
		return httperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to delete token, %v", err))
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		Domain:   ".inspirasisatu.com",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
	})
	return c.Status(fiber.StatusOK).JSON(
		response.WebResponse{
			Status:  fiber.StatusOK,
			Message: "Logout successful!",
		})
}
