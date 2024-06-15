package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string `gorm:"unique"`
	Password string
	Role     string
}

// LoginInput represents the required fields for login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
