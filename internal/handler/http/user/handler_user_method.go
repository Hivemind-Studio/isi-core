package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httphelper/response"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *UserHandler) GetStorageUnitDetail(c *fiber.Ctx) error {
	// Retrieve the "id" parameter from the request URL
	id, _ := strconv.Atoi(c.Params("id"))

	result, err := h.userService.GetTest(c, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(response.WebResponse{
				Status: id, Message: err.Error(),
			})
	}

	return c.Status(fiber.StatusOK).
		JSON(response.WebResponse{
			Status:  id,
			Message: "success",
			Data:    result,
		})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse the incoming JSON request body into a newUser object
	var newUser user.RegisterDTO
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(response.WebResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Invalid input",
			})
	}

	// Validate the newUser struct using the utility
	if err := utils.ValidateStruct(newUser); err != nil {
		// If validation fails, return a bad request response
		return c.Status(fiber.StatusBadRequest).
			JSON(response.WebResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Validation error",
				Data:    nil,
			})
	}

	// Call the userService to create the newUser
	result, err := h.userService.Create(c, &newUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(response.WebResponse{
				Status:  fiber.StatusBadRequest,
				Message: err.Error(),
			})
	}

	// Return success response with the created newUser data
	return c.Status(fiber.StatusCreated).
		JSON(response.WebResponse{
			Status:  fiber.StatusCreated,
			Message: "User created successfully",
			Data:    result,
		})
}
