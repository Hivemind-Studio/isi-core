package utils

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GenerateCookie creates a new cookie for the user after login
func GenerateCookie(c *fiber.Ctx, value string) {
	cookieName := os.Getenv("COOKIE_NAME")
	if cookieName == "" {
		cookieName = "default_cookie"
	}

	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour) // Cookie expires after 1 day
	cookie.HTTPOnly = true
	c.Cookie(cookie)
}
