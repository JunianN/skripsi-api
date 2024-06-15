package models

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	UserID         uint
	TranslatorID   uint // Reference to the User model for the assigned translator
	Title          string
	Description    string
	FilePath       string
	SourceLanguage string
	TargetLanguage string
	TranslatedPath string
	Status         string // Added to track translation status
}
