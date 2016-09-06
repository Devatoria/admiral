package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
)

// Namespace represents a namespace form model
type Namespace struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// getNamespaces returns all namespaces in database, ordered by name
func getNamespaces(c *gin.Context) {
	var namespaces []models.Namespace
	db.Instance().Order("name").Find(&namespaces)

	c.JSON(http.StatusOK, namespaces)
}

// getNamespace returns a namespace from the given ID
func getNamespace(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID format"})
		return
	}

	var namespace models.Namespace
	db.Instance().Where("id = ?", id).Find(&namespace)
	if namespace.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, namespace)
}

// putNamespace creates a namespace in database if it doesn't exist
func putNamespace(c *gin.Context) {
	var data Namespace
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	data.Name = strings.TrimSpace(data.Name)

	// Returns a conflict if namespace already exists
	var namespace models.Namespace
	db.Instance().Where("name = ?", data.Name).Find(&namespace)
	if namespace.ID != 0 {
		c.Status(http.StatusConflict)
		return
	}

	namespace.Name = data.Name
	db.Instance().Create(&namespace)

	c.JSON(http.StatusOK, namespace)
}

// deleteNamespace deletes the given namespace
func deleteNamespace(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID format"})
		return
	}

	db.Instance().Where("id = ?", id).Delete(&models.Namespace{})
	c.Status(http.StatusOK)
}

// patchNamespace updates the given
func patchNamespace(c *gin.Context) {
	// Check ID validity
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID format"})
		return
	}

	// Bind form to model
	var data Namespace
	err = c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	data.Name = strings.TrimSpace(data.Name)

	// Returns a conflict if namespace already exists
	var namespace models.Namespace
	db.Instance().Where("id = ?", id).Find(&namespace)
	if namespace.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	namespace.Name = data.Name
	db.Instance().Save(&namespace)

	c.JSON(http.StatusOK, namespace)
}

// getNamespaceImages returns a list of images associated to the given namespace
func getNamespaceImages(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID format"})
		return
	}

	var images []models.Image
	db.Instance().Where("namespace_id = ?", id).Find(&images)
	if len(images) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, images)
}
