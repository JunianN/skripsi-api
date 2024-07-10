package handlers

import (
	"path/filepath"
	"translation-app-backend/models"

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
		savePath := filepath.Join("uploads", "receipts", file.Filename)

		// Save the file to the server
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file: " + err.Error()})
		}

		// Store file path in the database
		document.PaymentReceiptFilePath = savePath
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
