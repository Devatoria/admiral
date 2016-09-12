package models

import (
	"github.com/jinzhu/gorm"
)

// Image represents an image stored in the registry
type Image struct {
	gorm.Model
	Name        string
	Namespace   Namespace `json:"-"`
	NamespaceID uint
}
