package models

import "gorm.io/gorm"

type Translation struct {
	gorm.Model
	DocumentID     uint
	TranslatorID   uint
	OriginalText   string
	TranslatedText string
	Status         string
}
