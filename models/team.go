package models

import (
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model
	Name  string `gorm:"not null;unique"`
	Users []User `gorm:"many2many:team_users;"`
}
