package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"time"
)

type AccessControlRule struct {
	Role          string
	AllowedMethod map[string][]string
}

type User struct {
	Name  string
	Email string
	Role  string
}

func validateUserRoles(accessControlRules map[string]AccessControlRule, apiEndpoint string, method string,
	user map[string]interface{}) (bool, error) {
	role, ok := user["createrole"].(string)
	if !ok || len(role) == 0 {
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

func JWTAuthMiddleware(accessControlRules map[string]AccessControlRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtToken := c.Get("Authorization")
		if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
			return fiber.ErrUnauthorized
		}

		jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

		claims, err := parseJWTToClaims(jwtToken)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Token")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
			}
		} else {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid expiration claim")
		}

		name := claims["name"].(string)
		email := claims["useremail"].(string)
		role := claims["createrole"].(string)

		user := map[string]interface{}{
			"name":       name,
			"useremail":  email,
			"createrole": role,
		}

		validRole, err := validateUserRoles(accessControlRules, c.Path(), c.Method(), user)
		if err != nil || !validRole {
			if err != nil {
				return fiber.ErrUnauthorized
			}
		}

		c.Locals("user_info", user)
		return c.Next()
	}
}

func parseJWTToClaims(tokenStr string) (jwt.MapClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fiber.ErrUnauthorized
}

func GenerateToken(user User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "defaultSecret"
	}

	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"useremail":  user.Email,
		"name":       user.Name,
		"createrole": user.Role,
		"exp":        expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
