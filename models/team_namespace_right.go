package models

import (
	_ "github.com/jinzhu/gorm"
)

// Team represents a team (group of users with specific rights)
type TeamNamespaceRight struct {
	Team        Team      `json:"-"`
	Namespace   Namespace `json:"-"`
	TeamID      uint      `gorm:"primary_key"`
	NamespaceID uint      `gorm:"primary_key"`
	Pull        bool
	Push        bool
}
