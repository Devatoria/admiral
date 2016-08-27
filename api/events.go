package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/docker/distribution/notifications"
	"github.com/gin-gonic/gin"
)

func getEvents(c *gin.Context) {
	var events []models.Event
	db.Instance().Find(&events)

	c.JSON(http.StatusOK, events)
}

func postEvents(c *gin.Context) {
	// Get data and unmarshall to envelop struct
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	var envelope notifications.Envelope
	err = json.Unmarshal(data, &envelope)
	if err != nil {
		panic(err)
	}

	// Create entities and persist
	for _, event := range envelope.Events {
		model := models.Event{
			EventID:          event.ID,
			Timestamp:        event.Timestamp,
			Action:           event.Action,
			TargetMediaType:  event.Target.MediaType,
			TargetSize:       event.Target.Size,
			TargetDigest:     event.Target.Digest.String(),
			TargetLength:     event.Target.Length,
			TargetRepository: event.Target.Repository,
			TargetURL:        event.Target.URL,
			TargetTag:        event.Target.Tag,
			RequestID:        event.Request.ID,
			RequestAddress:   event.Request.Addr,
			RequestHost:      event.Request.Host,
			RequestMethod:    event.Request.Method,
			RequestUserAgent: event.Request.UserAgent,
			SourceAddress:    event.Source.Addr,
			SourceInstanceID: event.Source.InstanceID,
		}

		db.Instance().Create(&model)
	}

	c.Status(http.StatusOK)
}
