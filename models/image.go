package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	Name        string
	Namespace   Namespace
	NamespaceID uint
}
