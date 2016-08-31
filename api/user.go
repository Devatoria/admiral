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

// getUsers returns users ordered by username
func getUsers(c *gin.Context) {
	nParam := c.DefaultQuery("n", "25")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var users []models.User
	db.Instance().Order("username").Limit(n).Find(&users)

	c.JSON(http.StatusOK, users)
}

// postUser creates a new user using the given data, and hashing the password using bcrypt
func postUser(c *gin.Context) {
	var data User
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Hash password before storing it in database
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		panic(err)
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

	user.Username = data.Username
	user.Password = string(hash)
	db.Instance().Create(&user)

	c.Status(http.StatusOK)
}
