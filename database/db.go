package database

import (
	"os"
	"translation-app-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	DB_URL := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(DB_URL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.Translation{}, &models.Notification{}, &models.Document{}, &models.Message{})
	return db, nil
}

// func AutoMigrate(db *gorm.DB) error {
// 	// Ensure all relevant models are listed here
// 	return db.AutoMigrate(&models.User{}, &models.Translation{}, &models.Notification{}, &models.Document{})
// }