package models

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	DocumentID uint   `gorm:"not null"`
	Message    string `gorm:"not null"`
	Read       bool   `gorm:"default:false"`
}
