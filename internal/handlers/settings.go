package handlers

import (
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdatePricePerWord(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			PricePerWord float64 `json:"price_per_word"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var settings models.Settings
		if err := db.First(&settings).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Settings not found"})
		}

		settings.PricePerWord = input.PricePerWord
		if err := db.Save(&settings).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update price"})
		}

		return c.JSON(fiber.Map{"message": "Price updated successfully", "price_per_word": settings.PricePerWord})
	}
} 