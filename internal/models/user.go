package models

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	RoleUser       = "user"
	RoleTranslator = "translator"
	RoleAdmin      = "admin"
)

type User struct {
	gorm.Model
	Username            string
	Email               string `gorm:"unique"`
	Password            string
	Role                string
	ProficientLanguages pq.StringArray `gorm:"type:text[]"`
	Categories          pq.StringArray `gorm:"type:text[];default:'{}'"`
	Ratings             []Rating       `gorm:"foreignKey:TranslatorID"`
	Status              string
}

// Validate validates user fields based on their role
func (u *User) Validate() error {
	if u.Role == RoleTranslator {
		for _, category := range u.Categories {
			switch category {
			case CategoryGeneral, CategoryEngineering, CategorySocialSciences:
				continue
			default:
				return errors.New("invalid category: must be one of 'general', 'engineering', or 'social sciences'")
			}
		}
	}
	return nil
}

// LoginInput represents the required fields for login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
