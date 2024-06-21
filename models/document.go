package models

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID             uint
	TranslatorID       uint // Reference to the User model for the assigned translator
	Title              string
	Description        string
	FilePath           string
	SourceLanguage     string
	TargetLanguage     string
	NumberOfPages      string
	TranslatedFilePath string // Path to the translated document
	Status             string // e.g., "Pending", "In Progress", "Completed"
	PaymentConfirmed   bool   // Field to check if the payment is confirmed
}
