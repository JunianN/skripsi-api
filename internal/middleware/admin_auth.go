package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AdminRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("userRole").(string)
        if !ok || role != "admin" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Access denied"})
        }
        return c.Next()
    }
}
