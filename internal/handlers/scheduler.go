package handlers

import (
	"log"
	"time"
	"translation-app-backend/internal/database"
	"translation-app-backend/internal/models"
)

func CheckAndDeclineUnconfirmedDocuments() {
	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return
	}
	
	var documents []models.Document
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	if err := db.Where("translator_approval_status = ? AND assignment_time < ?", "Pending", oneDayAgo).Find(&documents).Error; err != nil {
		log.Printf("Failed to fetch documents: %v", err)
		return
	}
	for _, document := range documents {
		document.TranslatorApprovalStatus = "Declined"
		if err := db.Save(&document).Error; err != nil {
			log.Printf("Failed to update document ID %d: %v", document.ID, err)
		} else {
			log.Printf("Document ID %d automatically declined due to no confirmation from translator", document.ID)
		}
	}
}
