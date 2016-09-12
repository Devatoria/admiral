package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/docker/distribution/notifications"
	"github.com/gin-gonic/gin"
)

func getEvents(c *gin.Context) {
	queryN := c.DefaultQuery("n", "25")
	n, err := strconv.Atoi(queryN)
	if err != nil {
		n = 25
	}

	var events []models.Event
	db.Instance().Order("created_at desc").Limit(n).Find(&events)

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
			ActorName:        event.Actor.Name,
		}

		db.Instance().Create(&model)

		// If event is a push, check if namespace and image exists or create them
		if event.Action == "push" && event.Target.Tag != "" {
			repSplit := strings.SplitN(event.Target.Repository, "/", 2)

			// Get namespace
			if len(repSplit) >= 2 {
				// Search for a namespace
				var namespace models.Namespace
				db.Instance().Where("name = ?", repSplit[0]).Find(&namespace)
				if namespace.ID == 0 {
					panic(fmt.Sprintf("Missing namespace: %s", repSplit[0]))
				}

				// Search for an image
				var image models.Image
				db.Instance().Where("name = ?", event.Target.Repository).Find(&image)
				if image.ID == 0 {
					image.Name = event.Target.Repository
					image.Namespace = namespace
					db.Instance().Create(&image)
				}

				// Search for a tag
				var tag models.Tag
				db.Instance().Where("name = ? AND image_id = ?", event.Target.Tag, image.ID).Find(&tag)
				if tag.ID == 0 {
					tag.Name = event.Target.Tag
					tag.Image = image
					db.Instance().Create(&tag)
				}
			}
		}
	}

	c.Status(http.StatusOK)
}
