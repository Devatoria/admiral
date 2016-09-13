package models

import (
	"github.com/jinzhu/gorm"
)

// User represents an... User... I think.
type User struct {
	gorm.Model
	Username     string `gorm:"not null;unique;index" json:"username" binding:"required"`
	Password     string `gorm:"-" json:"password" binding:"required"`
	PasswordHash string `gorm:"not null" json:"-"`
}
