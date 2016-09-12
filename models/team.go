package models

import (
	"github.com/jinzhu/gorm"
)

// Team represents a team (group of users with specific rights)
type Team struct {
	gorm.Model
	Name    string `gorm:"not null;unique"`
	Owner   User   `json:"-"`
	OwnerID uint
	Users   []User `gorm:"many2many:team_users;" json:"-"`
}
