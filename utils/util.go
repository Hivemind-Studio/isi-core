package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, dto interface{}) error {
	if err := c.BodyParser(dto); err != nil {
		return err
	}
	return nil
}

func SafeDereferenceString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToInterfaceSlice(intSlice []int64) []interface{} {
	result := make([]interface{}, len(intSlice))
	for i, v := range intSlice {
		result[i] = v
	}
	return result
}
