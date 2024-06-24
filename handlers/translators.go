package handlers

import (
    "github.com/gofiber/fiber/v2"
    "translation-app-backend/database"
    "translation-app-backend/models"
    "log"
)

func GetAssignedDocuments(c *fiber.Ctx) error {
    userID, _ := c.Locals("userID").(uint)
    userRole, _ := c.Locals("userRole").(string)

    if userRole != "translator" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
    }

    db, err := database.Connect()
    if err != nil {
        log.Printf("Database connection failed: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
    }

    var documents []models.Document
    if err := db.Where("translator_id = ?", userID).Find(&documents).Error; err != nil {
        log.Printf("Failed to fetch documents: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
    }

    return c.JSON(documents)
}
