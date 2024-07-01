package models

import (
    "gorm.io/gorm"
)

type Rating struct {
    gorm.Model
    UserID       uint   `gorm:"not null"`
    TranslatorID uint   `gorm:"not null"`
    DocumentID   uint   `gorm:"not null"`
    Rating       int    `gorm:"not null"`
    Comment      string `gorm:"type:text"`
}
