package handlers

import (
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Mail(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Message string `json:"message"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		mail := models.Mail{
			Name:    input.Name,
			Email:   input.Email,
			Message: input.Message,
		}

		if err := db.Create(&mail).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save mail form data"})
		}

		return c.JSON(fiber.Map{"message": "Your message has been sent successfully."})
	}
}

func GetMailSubmissions(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var mails []models.Mail
		if err := db.Find(&mails).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch mail submissions"})
		}

		return c.JSON(mails)
	}
}
