package database

import (
	"translation-app-backend/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Notification{}, &models.Document{}, &models.Discussion{}, &models.Rating{}, &models.Mail{},&models.Settings{})
}