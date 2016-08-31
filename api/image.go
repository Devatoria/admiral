package api

import (
	"net/http"
	"strconv"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

// getImages returns all images ordered by name
func getImages(c *gin.Context) {
	var images []models.Image
	db.Instance().Order("name").Find(&images)

	c.JSON(http.StatusOK, images)
}

// getImage returns the image associated to the given ID
func getImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID format"})
		return
	}

	var image models.Image
	db.Instance().Where("id = ?", id).Find(&image)
	if image.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, image)
}
