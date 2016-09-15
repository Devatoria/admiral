package models

import (
	"github.com/Devatoria/admiral/db"

	"github.com/jinzhu/gorm"
)

// Namespace represents an image namespace (kind of "group" of images)
type Namespace struct {
	gorm.Model
	Name    string `gorm:"not null;unique"`
	Owner   User   `json:"-"`
	OwnerID uint
	Images  []Image `json:"-"`
}

// GetNamespaceByName finds a namespace using the given name
func GetNamespaceByName(name string) Namespace {
	var namespace Namespace
	db.Instance().Preload("Images", func(gdb *gorm.DB) *gorm.DB {
		return gdb.Order("images.name")
	}).Preload("Images.Tags", func(gdb *gorm.DB) *gorm.DB {
		return gdb.Order("tags.name")
	}).Where("name = ?", name).Find(&namespace)

	return namespace
}
