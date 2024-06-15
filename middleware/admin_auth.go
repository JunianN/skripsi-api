package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AdminRequired checks if the user is an admin
func AdminRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Assuming user role is stored in Locals after authentication
        role, ok := c.Locals("userRole").(string)
        if !ok || role != "admin" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Access denied"})
        }
        return c.Next()
    }
}
