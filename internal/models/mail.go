package models

import (
    "gorm.io/gorm"
)

type Mail struct {
    gorm.Model
    Name    string `gorm:"not null"`
    Email   string `gorm:"not null"`
    Message string `gorm:"type:text;not null"`
}