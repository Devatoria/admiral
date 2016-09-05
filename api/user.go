package api

import (
	"net/http"
	"strconv"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/filters"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// getUsers returns users ordered by username and hide passwords
func getUsers(c *gin.Context) {
	nParam := c.DefaultQuery("n", "25")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var users []models.User
	db.Instance().Order("username").Limit(n).Find(&users)
	for i, _ := range users {
		users[i].Password = "[REDACTED]"
	}

	c.JSON(http.StatusOK, users)
}

// getUser returns the given user but hide its password
func getUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var user models.User
	db.Instance().Where("id = ?", id).Find(&user)
	if user.ID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	user.Password = "[REDACTED]"

	c.JSON(http.StatusOK, user)
}

// putUser creates a new user using the given data, and hashing the password using bcrypt
func putUser(c *gin.Context) {
	var data User
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Check that username does not already exist
	var user models.User
	db.Instance().Where("username = ?", data.Username).Find(&user)
	if user.ID != 0 {
		c.Status(http.StatusConflict)
		return
	}

	// Check username validity
	if err = filters.ValidateUsername(data.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check password validity
	if err = filters.ValidatePassword(data.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password before storing it in database
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		panic(err)
		return
	}

	user.Username = data.Username
	user.Password = string(hash)
	db.Instance().Create(&user)

	c.Status(http.StatusOK)
}
