package models

import (
	"gorm.io/gorm"
)

type Discussion struct {
	gorm.Model
	DocumentID uint
	UserID     uint
	Message    string
	UserRole   string // e.g., "user" or "admin"
}
