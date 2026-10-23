package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"strings"
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

func GenerateVerificationToken() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func IsOriginDashboard(origin string) bool {
	return strings.Contains(origin, "dashboard")
}

func IsOriginBackoffice(origin string) bool {
	return strings.Contains(origin, "backoffice")
}
