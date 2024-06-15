package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID  uint
	Message string
	Read    bool
}
