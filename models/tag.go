package models

import (
	"github.com/Devatoria/admiral/db"

	"github.com/jinzhu/gorm"
)

// Tag represents an image tag
type Tag struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Image   Image  `json:"-"`
	ImageID uint
}

// GetTagByName returns a tag using the given name and associated image
func GetTagByName(name string, image_id uint) Tag {
	var tag Tag
	db.Instance().Where("name = ? AND image_id = ?", name, image_id).Find(&tag)

	return tag
}
