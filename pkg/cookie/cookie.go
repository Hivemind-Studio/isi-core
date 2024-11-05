package cookie

import (
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func GenerateCookie(c *fiber.Ctx, user user.Cookie) {
	cookieName := os.Getenv("auth")
	if cookieName == "" {
		cookieName = "token"
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "defaultSecret"
	}

	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.RoleName,
		"exp":   expirationTime,
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create JWT token",
		})
		return
	}

	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.SameSite = "Strict"

	c.Cookie(cookie)
}
