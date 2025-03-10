package profile

import (
	"fmt"
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/s3"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "GetProfile"
	requestId, ok := c.Locals("request_id").(string)
	if !ok {
		requestId = "unknown"
	}

	logger.Print("info", requestId, module, functionName, "📌 [GetProfile] Request received", "")

	userSession := c.Locals("user")
	if userSession == nil {
		logger.Print("error", requestId, module, functionName,
			"user session is not found", userSession)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	user, ok := userSession.(*session.Session)
	if !ok {
		logger.Print("error", requestId, module, functionName, "Invalid session format", userSession)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid session data"})
	}

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Fetching profile for user: %+v", user), "")

	res, err := h.getProfileUser.Execute(c.Context(), user.ID, user.Role)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			fmt.Sprintf("Failed to retrieve profile for user ID=%s", user.ID), err.Error())
		return err
	}

	logger.Print("info", requestId, module, functionName,
		fmt.Sprintf("Successfully retrieved profile for user ID=%s", user.ID), "")

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile retrieved successfully",
		Data:    res,
	})
}

func (h *Handler) UpdateProfilePassword(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "UpdateProfilePassword"

	var requestBody authdto.UpdatePassword
	requestId := c.Locals("request_id").(string)

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	if err := validator.ValidatePassword(&requestBody.Password, &requestBody.ConfirmPassword); err != nil {
		return err
	}

	userSession, err := h.getUserSession(c)
	if err != nil {
		return err
	}

	err = h.updateProfilePassword.Execute(c.Context(), userSession.ID, requestBody.CurrentPassword, requestBody.Password)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "change password successfully",
	})
}

func (h *Handler) getUserSession(c *fiber.Ctx) (*session.Session, error) {
	userSession := c.Locals("user")
	if userSession == nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}
	user, ok := userSession.(*session.Session)
	if !ok {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid session data"})
	}
	return user, nil
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "UpdateProfile"
	requestId := c.Locals("request_id").(string)

	var requestBody dto.UpdateUserDTO

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	userSession, err := h.getUserSession(c)
	if err != nil {
		return err
	}

	res, err := h.updateProfile.Execute(c.Context(), userSession.ID, userSession.Role, requestBody)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile update successfully",
		Data:    res,
	})
}

func (h *Handler) UploadPhoto(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "UploadPhoto"
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName, "[UploadPhoto] Request received", "")

	userSession, err := h.getUserSession(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return httperror.New(fiber.StatusBadRequest, "Failed to retrieve file")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExts[ext] {
		return httperror.New(
			fiber.StatusBadRequest,
			"Invalid file type",
		)
	}

	tempFilePath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to save file")
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(tempFilePath)

	username := userSession.Name
	fileURL, err := s3.UploadFile(tempFilePath, file.Filename, username)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			"Failed to upload file to S3", err.Error())
		return httperror.New(fiber.StatusInternalServerError, "Failed to upload file")
	}

	err = h.updatePhoto.Execute(c.Context(), userSession.ID, fileURL)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile photo updated successfully",
		Data:    fileURL,
	})
}

func (h *Handler) DeletePhoto(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "DeletePhoto"
	requestId := c.Locals("request_id").(string)
	logger.Print("info", requestId, module, functionName, "[DeletePhoto] Request received", "")

	userSession, err := h.getUserSession(c)
	if err != nil {
		return err
	}

	err = h.deletePhoto.Execute(c.Context(), userSession.ID)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile photo updated successfully",
	})
}
