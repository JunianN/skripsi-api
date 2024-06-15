package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "translation-app-backend/models"
)

// ListTranslations - Lists all translations
func ListTranslations(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var translations []models.Translation
        db.Find(&translations)
        return c.JSON(translations)
    }
}

// AddTranslation - Adds a new translation
func AddTranslation(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        translation := new(models.Translation)
        if err := c.BodyParser(translation); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
        }

        db.Create(&translation)
        return c.JSON(translation)
    }
}
