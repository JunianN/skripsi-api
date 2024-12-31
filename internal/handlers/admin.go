package handlers

import (
	"log"
	"time"
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TranslatorWithRating struct {
	models.User
	AverageRating float64 `json:"average_rating"`
}

func RegisterAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if input.Email == "" || input.Password == "" || input.Username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "please provide all required fields"})
		}

		var exists models.User
		db.Where("email = ?", input.Email).First(&exists)
		if exists.ID != 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already in use"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot hash password"})
		}

		user := models.User{
			Username: input.Username,
			Email:    input.Email,
			Password: string(hashedPassword),
			Role:     "admin",
		}

		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create admin"})
		}

		return c.JSON(fiber.Map{"message": "New Admin registered successfully"})
	}
}

func GetAllDocuments(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var documents []models.Document
		if err := db.Order("created_at desc").Find(&documents).Error; err != nil {
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

		if document.FileContent == nil {
			log.Printf("Document file not found for document ID: %v", documentID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document file not found"})
		}

		// Set the appropriate headers
		c.Set("Content-Disposition", "attachment; filename="+document.FileName)
		c.Set("Content-Type", "application/octet-stream")

		return c.Send(document.FileContent)
	}
}

// ApproveDocument sets the approval status of a document to 'Approved'
func ApproveDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		document.ApprovalStatus = "Approved"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "Your document has been approved."
		if err := CreateNotification(document.UserID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Document approved successfully"})
	}
}

// RejectDocument sets the approval status of a document to 'Rejected'
func RejectDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		document.ApprovalStatus = "Rejected"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "Your document has been rejected."
		if err := CreateNotification(document.UserID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Document rejected successfully"})
	}
}

func AssignDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")
		var request struct {
			TranslatorID uint `json:"translator_id"`
		}
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		if document.ApprovalStatus != "Approved" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Document must be approved before assigning a translator"})
		}

		document.TranslatorID = request.TranslatorID
		document.TranslatorApprovalStatus = "Pending"
		document.AssignmentTime = time.Now() // Set the assignment time

		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "A document has been assigned to you."
		if err := CreateNotification(document.TranslatorID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Document assigned to translator"})
	}
}

func GetTranslators(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var translators []models.User
		if err := db.Where("role = ?", "translator").Find(&translators).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch translators"})
		}

		return c.JSON(translators)
	}
}

func DownloadTranslatedFile(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		if document.FileContent == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Translated document not found"})
		}

		// Set the appropriate headers
		c.Set("Content-Disposition", "attachment; filename="+document.TranslatedFileName)
		c.Set("Content-Type", "application/octet-stream")

		return c.Send(document.TranslatedFileContent)
	}
}

func ApproveTranslatedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		document.TranslatedApprovalStatus = "Approved"
		document.Status = "Finished"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "Your document has been translated."
		if err := CreateNotification(document.UserID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't create notification"})
		}

		var count int64
		db.Model(&models.Document{}).Where("translator_id = ? AND status = ?", document.TranslatorID, "Translating").Count(&count)

		if count == 0 {
			var translator models.User
			if err := db.First(&translator, document.TranslatorID).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Translator not found"})
			}

			translator.Status = "Available"
			if err := db.Save(&translator).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update translator status"})
			}
		}

		message2 := "Your translation has been approved by Admin."
		if err := CreateNotification(document.TranslatorID, document.ID, message2, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Translated document approved"})
	}
}

func RejectTranslatedDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		document.TranslatedApprovalStatus = "Rejected"
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "Your translation has been rejected by Admin."
		if err := CreateNotification(document.TranslatorID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Translated document rejected"})
	}
}

func DownloadPaymentReceipt(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		if document.PaymentReceiptContent == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment receipt not found"})
		}

		// Set the appropriate headers
		c.Set("Content-Disposition", "attachment; filename="+document.PaymentReceiptFileName)
		c.Set("Content-Type", "application/octet-stream")

		// Send the file content directly from the database
		return c.Send(document.PaymentReceiptContent)
	}
}

func ApprovePayment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var document models.Document
		if err := db.Where("id = ?", documentID).First(&document).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}

		document.PaymentConfirmed = true
		if err := db.Save(&document).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}

		message := "Your payment has been approved by Admin."
		if err := CreateNotification(document.UserID, document.ID, message, db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}

		return c.JSON(fiber.Map{"message": "Payment approved successfully"})
	}
}

func GetTranslatorsByLanguage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sourceLanguage := c.Query("source")
		targetLanguage := c.Query("target")
		category := c.Query("category")

		// Validate required parameters
		if sourceLanguage == "" || targetLanguage == "" || category == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "source language, target language, and category are required",
			})
		}

		// Validate category value
		switch category {
		case models.CategoryGeneral, models.CategoryEngineering, models.CategorySocialSciences:
			// Valid category, proceed
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid category: must be one of 'general', 'engineering', or 'social sciences'",
			})
		}

		query := `
			SELECT users.*, COALESCE(AVG(ratings.rating), 0) as average_rating
			FROM users
			LEFT JOIN ratings ON users.id = ratings.translator_id
			WHERE users.role = ?
			AND ARRAY[?] <@ users.proficient_languages 
			AND ARRAY[?] <@ users.proficient_languages
			AND ARRAY[?] <@ users.categories
			GROUP BY users.id
			ORDER BY average_rating DESC, users.status ASC`

		var translators []TranslatorWithRating
		if err := db.Raw(query, "translator", sourceLanguage, targetLanguage, category).Scan(&translators).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch translators"})
		}

		return c.JSON(translators)
	}
}

func DeleteTranslator(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		translatorID := c.Params("id")

		// Check if translator exists and is actually a translator
		var translator models.User
		if err := db.Where("id = ? AND role = ?", translatorID, "translator").First(&translator).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Translator not found"})
		}

		// Check if translator has any ongoing translations
		var ongoingDocuments int64
		if err := db.Model(&models.Document{}).Where("translator_id = ? AND status = ?", translatorID, "Translating").Count(&ongoingDocuments).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check ongoing translations"})
		}

		if ongoingDocuments > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot delete translator with ongoing translations"})
		}

		// Delete the translator
		if err := db.Delete(&translator).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete translator"})
		}

		return c.JSON(fiber.Map{"message": "Translator deleted successfully"})
	}
}
