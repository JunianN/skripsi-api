package handlers

import (
	"strconv"
	"translation-app-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetDiscussions retrieves all discussion messages for a specific document
func GetDiscussions(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var discussions []models.Discussion
		if err := db.Where("document_id = ?", documentID).Find(&discussions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch discussions"})
		}

		return c.JSON(discussions)
	}
}

// PostDiscussion adds a new discussion message for a specific document
func PostDiscussion(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(float64)
		userRole, _ := c.Locals("userRole").(string)
		documentIDStr := c.Params("id")
		documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid document ID"})
		}

		var input struct {
			Message string `json:"message"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		discussion := models.Discussion{
			DocumentID: uint(documentID),
			UserID:     uint(userID),
			Message:    input.Message,
			UserRole:   userRole,
		}

		if err := db.Create(&discussion).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create discussion"})
		}

		return c.JSON(discussion)
	}
}
