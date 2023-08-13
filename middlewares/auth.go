package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		c.Locals("token", token)
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUnauthorized)

}
