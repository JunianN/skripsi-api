package models

import (
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	PricePerWord float64 `gorm:"not null"`
} 