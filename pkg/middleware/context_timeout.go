package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}
