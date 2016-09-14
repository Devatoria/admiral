package models

import (
	"github.com/Devatoria/admiral/db"

	"github.com/jinzhu/gorm"
)

// Image represents an image stored in the registry
type Image struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	Namespace   Namespace `json:"-"`
	NamespaceID uint
	Tags        []Tag
}

// GetImageByName returns an image using the given name
func GetImageByName(name string) Image {
	var image Image
	db.Instance().Where("name = ?", name).Find(&image)

	return image
}
