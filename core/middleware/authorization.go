package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Authorizes() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHeader := c.Get("Authorization")
		var token string
		if tokenHeader != "" && strings.HasPrefix(tokenHeader, "Bearer ") {
			token = strings.TrimPrefix(tokenHeader, "Bearer ")
		} else {
			token = c.Cookies("access_token")        // read cookie set by SetCookies
		}
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}
		jwtWrapper := JwtWrapper{
			SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
			Issuer:          "authService",
			ExpirationHours: 24,
		}
		// `token` now contains the raw JWT (either trimmed from "Bearer " header or read from cookie)
		_, err := jwtWrapper.ValidateToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		return c.Next()
	}
}

func SetCookies(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}

func CORS(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Next()
	})
}
