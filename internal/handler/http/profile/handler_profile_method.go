package profile

import (
	"fmt"
	authdto "github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/s3"
	"github.com/Hivemind-Studio/isi-core/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	user, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	res, err := h.getProfileUser.Execute(c.Context(), user.ID, user.Role)

	if err != nil {
		return err
	}

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

	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if err := validator.ValidatePassword(&requestBody.Password, &requestBody.ConfirmPassword); err != nil {
		return err
	}

	err = h.updateProfilePassword.Execute(c.Context(), claims.ID, requestBody.CurrentPassword, requestBody.Password)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "change password successfully",
	})
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

	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	res, err := h.updateProfile.Execute(c.Context(), claims.ID, requestBody.Name, requestBody.Address,
		requestBody.Gender, requestBody.PhoneNumber, requestBody.Occupation)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile update successfully",
		Data:    res,
	})
}

func (h *Handler) UpdateProfileCoach(c *fiber.Ctx) error {
	module := "Profile Handler"
	functionName := "UpdateProfileCoach"
	requestId := c.Locals("request_id").(string)

	var requestBody dto.UpdateCoachDTO

	if err := c.BodyParser(&requestBody); err != nil {
		logger.Print("error", requestId, module, functionName,
			"Invalid input", string(c.Body()))
		return httperror.New(fiber.StatusBadRequest, "Invalid input")
	}

	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	res, err := h.updateProfileCoach.Execute(c.Context(), claims.ID, requestBody.Name, requestBody.Address,
		requestBody.Gender, requestBody.PhoneNumber, requestBody.DateOfBirth, requestBody.Title, requestBody.Bio, requestBody.Expertise)

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

	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
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

	username := claims.Name
	fileURL, err := s3.UploadFile(tempFilePath, file.Filename, username)
	if err != nil {
		logger.Print("error", requestId, module, functionName,
			"Failed to upload file to S3", err.Error())
		return httperror.New(fiber.StatusInternalServerError, "Failed to upload file")
	}

	err = h.updatePhoto.Execute(c.Context(), claims.ID, fileURL)
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
	jwtToken := c.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return fiber.ErrUnauthorized
	}

	claims, err := middleware.ExtractJWTPayload(jwtToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	err = h.deletePhoto.Execute(c.Context(), claims.ID)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.WebResponse{
		Status:  fiber.StatusOK,
		Message: "Profile photo updated successfully",
	})
}
