package models

import (
	"github.com/jinzhu/gorm"
)

// Tag represents an image tag
type Tag struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Image   Image  `json:"-"`
	ImageID uint
}
