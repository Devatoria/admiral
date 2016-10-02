package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getLogin only returns 200 if the user is able to login, 401 otherwise
func getLogin(c *gin.Context) {
	c.Status(http.StatusOK)
}
