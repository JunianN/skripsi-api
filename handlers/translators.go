package handlers

import (
	"path/filepath"
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAssignedDocuments(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		var documents []models.Document
		if err := db.Where("translator_id = ?", userID).Find(&documents).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
		}

		return c.JSON(documents)
	}
}

func GetAssignedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
		}

		return c.JSON(document)
	}
}

func DownloadAssignedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
		}

		if document.FilePath == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File path not found"})
		}

		filename := filepath.Base(document.FilePath)
		c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		err := c.SendFile(document.FilePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send file"})
		}

		return nil
	}
}
