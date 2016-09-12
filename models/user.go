package models

import (
	"github.com/jinzhu/gorm"
)

// User represents a... User... I think.
type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null" json:"-"`
	IsAdmin  bool   `gorm:not null"`
	Teams    []Team `gorm:"many2many:team_users;" json:"-"`
}
