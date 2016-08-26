package models

import (
	"github.com/jinzhu/gorm"
)

type EventSource struct {
	gorm.Model
	Adress     string
	InstanceID string
}
