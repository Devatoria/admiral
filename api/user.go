package api

import (
	"net/http"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/filters"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func putUser(c *gin.Context) {
	var err error

	// Bind model
	var data models.User
	if c.BindJSON(&data) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Sanitize
	data.Username, err = filters.SanitizeUsername(data.Username)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// None of these should exist: user and namespace
	if db.Exists(db.Instance(), "username", data.Username, &models.User{}) {
		c.Status(http.StatusConflict)
		return
	}
	if db.Exists(db.Instance(), "name", data.Username, &models.Namespace{}) {
		c.Status(http.StatusConflict)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	data.PasswordHash = string(hash)

	// Create user and namespace
	db.Instance().Create(&data)
	db.Instance().Create(&models.Namespace{
		Name:  data.Username,
		Owner: data,
	})

	c.Status(http.StatusOK)
}
