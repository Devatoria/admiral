package models

import (
	"github.com/jinzhu/gorm"
)

type EventRequest struct {
	gorm.Model
	ReqID     string
	Address   string
	Host      string
	Method    string
	UserAgent string
}
