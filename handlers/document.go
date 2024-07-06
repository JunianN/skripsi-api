package handlers

import (
	"path/filepath"
	"strconv"
	"translation-app-backend/database"
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetDocument retrieves the details of a specific document by ID
func GetDocument(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	documentID := c.Params("id")

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	var document models.Document
	if err := db.Where("id = ? AND user_id = ?", documentID, userID).First(&document).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	return c.JSON(document)
}

// GetDocuments returns a list of documents for the authenticated user
func GetDocuments(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID") // Assuming userID is stored in Locals after authentication

		var documents []models.Document
		result := db.Where("user_id = ?", userID).Order("created_at desc").Find(&documents)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve documents"})
		}

		return c.JSON(documents)
	}
}

// UploadDocument handles the uploading of files along with additional data
func UploadDocument(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if userID is set in locals (set by middleware)
		userID := c.Locals("userID").(float64)

		// Parse the form
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form: " + err.Error()})
		}

		// Extract file from the posted form
		files := form.File["document"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
		}

		file := files[0]
		savePath := filepath.Join("uploads", file.Filename)

		// Save the file to the server
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file: " + err.Error()})
		}

		// Extract other form fields
		title := form.Value["title"][0]
		description := form.Value["description"][0]
		sourceLanguage := form.Value["sourceLanguage"][0]
		targetLanguage := form.Value["targetLanguage"][0]
		numberOfPages := form.Value["numberOfPages"][0]

		numberOfPagesInt, err := strconv.Atoi(numberOfPages)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to convert string: " + err.Error()})
		}

		if numberOfPagesInt <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Number of pages must be a positive integer"})
		}

		doc := models.Document{
			UserID:         uint(userID),
			Title:          title,
			Description:    description,
			FilePath:       savePath,
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
			NumberOfPages:  numberOfPages,
			Status:         "Pending", // Default status set when uploading a new document
		}

		if err := db.Create(&doc).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create document" + err.Error()})
		}

		return c.JSON(fiber.Map{"message": "File uploaded successfully", "data": doc})
	}
}

// UpdateDocumentStatus allows a translator to accept or decline a document assignment
func UpdateDocumentStatus(c *fiber.Ctx) error {
	documentID := c.Params("id")
	userID := c.Locals("userID") // Retrieved from authenticated session

	var input struct {
		Status string `json:"status"` // Accepts "Accepted" or "Declined"
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validate status input
	if input.Status != "Accepted" && input.Status != "Declined" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status provided"})
	}

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	// Fetch the document to ensure it's assigned to the current user
	var document models.Document
	if err := db.Where("id = ? AND translator_id = ?", documentID, userID).First(&document).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found or not assigned to you"})
	}

	// Update the status
	document.Status = input.Status
	if err := db.Save(&document).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document status"})
	}

	return c.JSON(fiber.Map{"message": "Document status updated successfully", "status": document.Status})
}

// GetTranslatorDocuments retrieves documents assigned to the logged-in translator
func GetTranslatorDocuments(c *fiber.Ctx) error {
	translatorID := c.Locals("userID") // Assuming userID is the ID of the translator

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	var documents []models.Document
	db = db.Where("translator_id = ?", translatorID).Find(&documents)

	if db.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(documents)
}

// DownloadTranslatedDocument handles the downloading of the translated document
func DownloadTranslatedDocument(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	documentID := c.Params("id")

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	var document models.Document
	if err := db.Where("id = ? AND user_id = ?", documentID, userID).First(&document).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	if !document.PaymentConfirmed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Payment not confirmed"})
	}

	return c.SendFile(document.TranslatedFilePath)
}
