package middleware

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks if the user is authenticated by verifying the cookie
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get cookie name from environment variable
		cookieName := os.Getenv("COOKIE_NAME")
		if cookieName == "" {
			cookieName = "default_cookie"
		}

		// Retrieve the cookie value
		cookie := c.Cookies(cookieName)
		if cookie == "" {
			return c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
		}

		// Here, you can add additional validation logic for the cookie,
		// such as checking a token against a database or verifying its signature.

		// If the cookie is valid, proceed to the next handler
		return c.Next()
	}
}
