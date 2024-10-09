package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ParseBody is a generic function to parse request body into the provided struct reference
func ParseBody(c *fiber.Ctx, dto interface{}) error {
	if err := c.BodyParser(dto); err != nil {
		return err
	}
	return nil
}
