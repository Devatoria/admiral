package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getVersion(c *gin.Context) {
	c.String(http.StatusOK, "0.1.0")
}
