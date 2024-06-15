package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	DocumentID uint
	UserID     uint
	Text       string
}
