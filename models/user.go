package models

import (
	"github.com/jinzhu/gorm"
)

// User represents a... User... I think.
type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:not null"`
}
