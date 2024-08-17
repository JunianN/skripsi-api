package handlers

import (
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FetchNotifications(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		var notifications []models.Notification
		if err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch notifications"})
		}

		return c.JSON(notifications)
	}
}

func CreateNotification(userID uint, documentID uint, message string, db *gorm.DB) error {
		notification := models.Notification{
			UserID:  userID,
			DocumentID: documentID,
			Message: message,
		}

		if err := db.Create(&notification).Error; err != nil {
			return err
		}

		return nil
	}

func MarkNotificationsAsRead(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		if err := db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("read", true).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark notifications as read"})
		}

		return c.JSON(fiber.Map{"message": "All notifications marked as read"})
	}
}
