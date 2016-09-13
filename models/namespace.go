package models

import (
	"github.com/jinzhu/gorm"
)

// Namespace represents an image namespace (kind of "group" of images)
type Namespace struct {
	gorm.Model
	Name    string `gorm:"not null;unique"`
	Owner   User   `json:"-"`
	OwnerID uint
}
