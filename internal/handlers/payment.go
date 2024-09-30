package handlers

import (
	"io"
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UploadPaymentReceipt(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ? AND user_id = ?", documentID, userID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		// Parse the form
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form: " + err.Error()})
		}

		// Extract file from the posted form
		files := form.File["receipt"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
		}

		file := files[0]

		// Read the file content
		fileContent, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file: " + err.Error()})
		}
		defer fileContent.Close()

		// Read the file data
		fileData, err := io.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file: " + err.Error()})
		}

		// Store file content in the database
		document.PaymentReceiptContent = fileData
		document.PaymentReceiptFileName = file.Filename
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "User has uploaded the payment receipt."
		if err := CreateNotification(2, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Payment receipt uploaded successfully"})
	}
}
