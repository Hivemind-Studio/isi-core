package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"strconv"
	"strings"
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

	log.Info().
		Str("request_id", requestID).
		Str("path", path).
		Interface("headers", headers).
		Str("body", body).
		Msg("Incoming request details")
}

func BodyLimit(limitMB string) fiber.Handler {
	limitStr := strings.TrimSuffix(strings.TrimSpace(limitMB), "MB")
	megabytes, err := strconv.Atoi(limitStr)
	if err != nil {
		panic("invalid body limit parameter: must be in format '3MB'")
	}

	bytes := megabytes * 1024 * 1024

	return func(c *fiber.Ctx) error {
		length := c.Request().Header.ContentLength()

		if length > bytes {
			return fiber.ErrRequestEntityTooLarge
		}

		if bodyStream := c.Request().BodyStream(); bodyStream != nil {
			c.Request().SetBodyStream(
				&io.LimitedReader{
					R: bodyStream,
					N: int64(bytes),
				},
				length,
			)
		}

		return c.Next()
	}
}
