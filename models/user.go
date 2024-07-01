package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string
	Email               string `gorm:"unique"`
	Password            string
	Role                string         // e.g., "admin", "translator", "user"
	ProficientLanguages pq.StringArray `gorm:"type:text[]"`
	Ratings             []Rating       `gorm:"foreignKey:TranslatorID"`
	Status              string         // e.g., "available", "working"
}

// LoginInput represents the required fields for login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
