package googlelogin

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"time"
)

func generateStateOauthCookie() (string, fiber.Cookie) {
	expiration := time.Now().Add(24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := fiber.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HTTPOnly: true,
		Secure:   true,
	}

	return state, cookie
}

func (uc *UseCase) Execute(c *fiber.Ctx) string {
	state, cookie := generateStateOauthCookie()

	c.Cookie(&cookie)

	url := uc.Oauth2.AuthCodeURL(state)

	return url
}
