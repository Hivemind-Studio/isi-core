package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// User represents the structure of the user information in the cookie
type User struct {
	Role   string `json:"role"`
	UserID string `json:"userID"`
}

// RoleMiddleware Middleware to enforce role-based access
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the role from the "ISI" cookie
		cookieValue := c.Cookies("ISI")

		if cookieValue == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		var user User
		// Deserialize the JSON string from the cookie
		if err := json.Unmarshal([]byte(cookieValue), &user); err != nil {
			fmt.Println("Error parsing cookie:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		fmt.Println("User Role:", user.Role)

		// Check if the user's role is in the allowed roles
		for _, role := range allowedRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}
}
