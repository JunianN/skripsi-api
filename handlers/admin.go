package handlers

import (
	"log"
	"path/filepath"
	"time"
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

		if document.TranslatedFilePath == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Translated document not found"})
		}

		err := c.SendFile(document.TranslatedFilePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send file"})
		}

		return nil
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

		if document.PaymentReceiptFilePath == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment receipt not found"})
		}

		filename := filepath.Base(document.PaymentReceiptFilePath)
		c.Set("Content-Disposition", "attachment; filename="+filename)

		err := c.SendFile(document.PaymentReceiptFilePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send file"})
		}

		return nil
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

		return c.JSON(fiber.Map{"message": "Payment approved successfully"})
	}
}

func GetTranslatorsByLanguage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sourceLanguage := c.Query("source")
		targetLanguage := c.Query("target")

		var translators []models.User
		if err := db.Where("role = ? AND proficient_languages @> ARRAY[?] AND proficient_languages @> ARRAY[?]", "translator", sourceLanguage, targetLanguage).Find(&translators).Error; err != nil {
			log.Printf("Failed to fetch translators: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch translators"})
		}

		return c.JSON(translators)
	}
}
