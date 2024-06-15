package handlers

import (
	"strconv"
	"translation-app-backend/database"
	"translation-app-backend/models"

	"github.com/gofiber/fiber/v2"
)

func AddMessage(c *fiber.Ctx) error {
	documentID := c.Params("id") // Get document ID from URL
	userID, _ := c.Locals("userID").(float64)

	var input struct {
		Text string `json:"text"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	docID, err := strconv.ParseUint(documentID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid document ID"})
	}

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	message := models.Message{
		DocumentID: uint(docID),
		UserID:     uint(userID),
		Text:       input.Text,
	}

	db.Create(&message)

	return c.JSON(fiber.Map{"message": "Message added successfully", "data": message})
}

func GetMessages(c *fiber.Ctx) error {
	documentID := c.Params("id")

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection failed"})
	}

	var messages []models.Message
	db.Where("document_id = ?", documentID).Find(&messages)

	return c.JSON(messages)
}
