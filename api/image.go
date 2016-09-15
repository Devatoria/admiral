package api

import (
	"fmt"
	"net/http"

	"github.com/Devatoria/admiral/auth"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

// getImages returns the images contained in the user namespace
func getImages(c *gin.Context) {
	// Get user
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Get user namespace
	namespace := models.GetNamespaceByName(user.Username)
	if namespace.ID == 0 {
		panic(fmt.Sprintf("User %s has no namespace", user.Username))
	}

	// Get images
	c.JSON(http.StatusOK, namespace.Images)
}
