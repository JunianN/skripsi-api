package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string
	Email               string 			`gorm:"unique"`
	Password            string
	Role                string         
	ProficientLanguages pq.StringArray `gorm:"type:text[]"`
	Ratings             []Rating       `gorm:"foreignKey:TranslatorID"`
	Status              string         
}

// LoginInput represents the required fields for login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
