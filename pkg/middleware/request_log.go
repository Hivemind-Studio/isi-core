package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

func RequestIdMiddleware(c *fiber.Ctx) error {
	requestID := c.Get("request_id")

	if requestID == "" {
		requestID = uuid.New().String()
	}

	c.Locals("request_id", requestID)

	logRequestDetails(c, requestID)

	return c.Next()
}

func logRequestDetails(c *fiber.Ctx, requestID string) {
	path := c.Path()

	headers := c.GetReqHeaders()

	var body string
	if c.Body() != nil {
		body = string(c.Body())
	} else {
		body = "no body"
	}

	log.Printf("Request ID: %s | Path: %s | Headers: %v | Body: %s", requestID, path, headers, body)
}
