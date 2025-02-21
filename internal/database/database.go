package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	DB_URL := os.Getenv("DB_URL")

	pgConfig := postgres.Config{
		DSN: DB_URL,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		PrepareStmt: false, 
	})
	if err != nil {
		return nil, err
	}

	Migrate(db)

	return db, nil
}