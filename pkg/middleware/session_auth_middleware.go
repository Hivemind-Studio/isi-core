package middleware

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/constant/rediskey"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strings"
)

func SessionAuthMiddleware(sessionManager *session.SessionManager, accessControlRules map[string]AccessControlRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		environment := os.Getenv("ENVIRONMENT")
		origin := c.Get("Origin")
		appOrigin := c.Get(constant.APP_ORIGIN_HEADER)

		log.Printf("SessionAuthMiddleware: ENVIRONMENT=%s, Origin=%s, X-App-Origin=%s", environment, origin, appOrigin)

		var cookieName string

		if environment == constant.PRODUCTION {
			if utils.IsOriginDashboard(origin) || appOrigin == constant.DASHBOARD {
				cookieName = constant.TOKEN_ACCESS_DASHBOARD
			} else if utils.IsOriginBackoffice(origin) || appOrigin == constant.BACKOFFICE {
				cookieName = constant.TOKEN_ACCESS_BACKOFFICE
			} else {
				log.Println("Unauthorized access: Invalid token origin in production")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token origin"})
			}
		} else {
			if appOrigin == constant.DASHBOARD {
				cookieName = constant.TOKEN_ACCESS_DASHBOARD
			} else if appOrigin == constant.BACKOFFICE {
				cookieName = constant.TOKEN_ACCESS_BACKOFFICE
			} else {
				log.Println("Unauthorized access: Invalid token origin in non-production")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token origin"})
			}
		}

		token := c.Cookies(cookieName)
		sessionID := rediskey.SESSION_PREFIX_KEY + token

		log.Printf("SessionAuthMiddleware: Selected cookie=%s, SessionID=%s", cookieName, sessionID)

		if token == "" {
			log.Println("Unauthorized access: No session token found in cookies")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No session"})
		}

		ctx := context.Background()
		userSession, err := sessionManager.GetSession(ctx, sessionID)
		if err != nil {
			log.Printf("Unauthorized access: Invalid session for sessionID=%s, Error=%v", sessionID, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session"})
		}

		validRole, err := validateUserRoles(accessControlRules, c.Path(), c.Method(), userSession)
		if err != nil {
			log.Printf("Unauthorized access: Role validation failed for sessionID=%s, Error=%v", sessionID, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid path for role"})
		}
		if !validRole {
			log.Printf("Unauthorized access: User sessionID=%s does not have a valid role for %s %s", sessionID, c.Method(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid path for role"})
		}

		log.Printf("SessionAuthMiddleware: User authorized, sessionID=%s, Path=%s, Method=%s", sessionID, c.Path(), c.Method())

		c.Locals("user", userSession)
		c.Locals("token", token)
		return c.Next()
	}
}

type AccessControlRule struct {
	Role          string
	AllowedMethod map[string][]string
}

type User struct {
	ID    int64
	Name  string
	Email string
	Role  string
}

func validateUserRoles(accessControlRules map[string]AccessControlRule, apiEndpoint string, method string,
	user *session.Session) (bool, error) {
	role := user.Role
	if len(role) == 0 {
		return false, fmt.Errorf("empty or invalid roles")
	}

	for _, rule := range accessControlRules {
		if role == rule.Role {
			for path, allowedMethods := range rule.AllowedMethod {
				if strings.HasPrefix(apiEndpoint, path) && sliceContainsString(allowedMethods, method) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func sliceContainsString(slice []string, target string) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	for _, val := range slice {
		if val == target {
			return true
		}
	}
	return false
}
