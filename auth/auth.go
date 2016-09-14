package auth

import (
	"errors"
	"net/http"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate returns an error if the given request basic auth is unable to authenticate the user, nil otherwise
func Authenticate(req *http.Request) (models.User, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return models.User{}, errors.New("Unable to get request basic auth")
	}

	var user models.User
	db.Instance().Where("username = ?", username).Find(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	return user, err
}

// GetCurrentUser returns the current user entity if authenticated
func GetCurrentUser(c *gin.Context) (models.User, error) {
	user, ok := c.Keys["user"].(models.User)
	if !ok {
		return models.User{}, errors.New("Unable to retrieve current user")
	}

	return user, nil
}
