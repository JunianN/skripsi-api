package handlers

import (
	"io"
	"translation-app-backend/internal/models"

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

		if document.FileContent == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document file not found"})
		}

		c.Set("Content-Disposition", "attachment; filename="+document.FileName)
		c.Set("Content-Type", "application/octet-stream")

		return c.Send(document.FileContent)
	}
}

func ApproveAssignedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
		}

		document.TranslatorApprovalStatus = "Accepted"
		document.Status = "Translating"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		var translator models.User
		if err := db.First(&translator, userID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Translator not found"})
		}

		translator.Status = "Working"
		if err := db.Save(&translator).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update translator status"})
		}

		
		message := "Your document is being translated."
		if err := CreateNotification(document.UserID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		message2 := "A translator has accepted to translate a document."
		if err := CreateNotification(2, document.ID, message2, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		
		return c.JSON(fiber.Map{"message": "Document accepted successfully"})
	}
}

func DeclineAssignedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
		}

		document.TranslatorApprovalStatus = "Declined"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message2 := "A translator has refused to translate a document."
		if err := CreateNotification(2, document.ID, message2, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Document declined successfully"})
	}
}

func UploadTranslatedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
		}

		file, err := c.FormFile("translated_document")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
		}

		fileContent, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file: " + err.Error()})
		}
		defer fileContent.Close()

		/// Read the file data
		fileData, err := io.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file: " + err.Error()})
		}

		document.TranslatedFileContent = fileData
		document.TranslatedFileName = file.Filename
		document.TranslatedApprovalStatus = "Pending"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "A translator has submited translated document."
		if err := CreateNotification(2, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Translated document uploaded successfully"})
	}
}
