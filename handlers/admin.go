package handlers

import (
	"log"
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllDocuments(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var documents []models.Document
		if err := db.Find(&documents).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
		}

		return c.JSON(documents)
	}
}

func GetDocumentDetails(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		return c.JSON(document)
	}
}

func DownloadUserDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			log.Printf("Document not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		if document.FilePath == "" {
			log.Printf("Document file not found for document ID: %v", documentID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document file not found"})
		}

		err := c.SendFile(document.FilePath)
		if err != nil {
			log.Printf("Error sending file: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send file"})
		}

		return nil
	}
}
