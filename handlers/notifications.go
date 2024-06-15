package handlers

import (
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SendNotification - Sends a notification to a user
func SendNotification(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notification := new(models.Notification)
		if err := c.BodyParser(notification); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		db.Create(&notification)
		return c.JSON(notification)
	}
}
