package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

// GenerateCookie creates a new cookie for the user after login
func GenerateCookie(c *fiber.Ctx, userID string, role string) {
	cookieName := os.Getenv("COOKIE_NAME")
	if cookieName == "" {
		cookieName = "token"
	}

	// Create the cookie value by appending the role to the userID
	value, err := json.Marshal(map[string]string{
		"userID": userID,
		"role":   role,
	})
	if err != nil {
		// Handle JSON marshaling error
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create cookie value",
		})
		return
	}

	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = string(value)                    // Store the JSON string as the cookie value
	cookie.Expires = time.Now().Add(24 * time.Hour) // Cookie expires after 1 day
	cookie.HTTPOnly = true                          // Prevent access to the cookie from JavaScript
	cookie.Secure = true                            // Ensure this is true in production (HTTPS)
	cookie.SameSite = "Strict"                      // CSRF protection

	c.Cookie(cookie)
}
