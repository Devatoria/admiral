package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(address string, port int) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/version", getVersion)

	// Registry events notification endpoint
	r.POST("/events", postEvents)

	r.Run(fmt.Sprintf("%s:%d", address, port))
}
