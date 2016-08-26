package models

import (
	"github.com/jinzhu/gorm"
)

type EventTarget struct {
	gorm.Model
	MediaType  string
	Size       int
	Digest     string
	Length     int
	Repository string
	Url        string
	Tag        string
}
