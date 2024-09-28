package database

import (
	"os"
	"translation-app-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	DB_URL := os.Getenv("DB_URL")

	// Create a new PostgreSQL driver configuration
	pgConfig := postgres.Config{
		DSN: DB_URL,
		PreferSimpleProtocol: true, // Disable implicit prepared statement usage
	}

	// Open the database connection with the modified configuration
	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		PrepareStmt: false, // Disable prepared statement cache
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{}, &models.Notification{}, &models.Document{}, &models.Discussion{}, &models.Rating{}, &models.Mail{})
	return db, nil
}