package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	EventID          string
	Timestamp        time.Time
	Action           string
	TargetMediaType  string
	TargetSize       int64
	TargetDigest     string
	TargetLength     int64
	TargetRepository string
	TargetURL        string
	TargetTag        string
	RequestID        string
	RequestAddress   string
	RequestHost      string
	RequestMethod    string
	RequestUserAgent string
	SourceAddress    string
	SourceInstanceID string
}
