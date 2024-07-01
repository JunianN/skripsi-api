package handlers

import (
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubmitRating(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			TranslatorID uint   `json:"translator_id"`
			DocumentID   uint   `json:"document_id"`
			Rating       int    `json:"rating"`
			Comment      string `json:"comment"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		userID := c.Locals("userID").(float64)

		rating := models.Rating{
			UserID:       uint(userID),
			TranslatorID: input.TranslatorID,
			DocumentID:   input.DocumentID,
			Rating:       input.Rating,
			Comment:      input.Comment,
		}

		if err := db.Create(&rating).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create rating"})
		}

		return c.JSON(fiber.Map{"message": "rating submitted successfully"})
	}
}

func GetRatings(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")

		var ratings []models.Rating
		if err := db.Where("document_id = ?", documentID).Find(&ratings).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch ratings"})
		}

		return c.JSON(ratings)
	}
}

func GetTranslatorAverageRating(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		translatorID := c.Params("id")

		var avgRating float64
		if err := db.Table("ratings").Where("translator_id = ?", translatorID).Select("AVG(rating)").Row().Scan(&avgRating); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch average rating"})
		}

		return c.JSON(fiber.Map{"average_rating": avgRating})
	}
}
