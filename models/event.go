package models

import (
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	EventID   string
	Timestamp string
	Action    string
	Target    EventTarget
	Request   EventRequest
	Source    EventSource
}
