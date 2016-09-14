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
}

// GetNamespaceByName finds a namespace using the given name
func GetNamespaceByName(name string) Namespace {
	var namespace Namespace
	db.Instance().Where("name = ?", name).Find(&namespace)

	return namespace
}
