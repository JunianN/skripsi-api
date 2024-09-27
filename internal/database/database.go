package database

import (
	"os"
	"translation-app-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	DB_URL := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(DB_URL), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, err
	}
	db.Exec("DEALLOCATE ALL")
	db.AutoMigrate(&models.User{}, &models.Notification{}, &models.Document{}, &models.Discussion{}, &models.Rating{}, &models.Mail{})
	return db, nil
}