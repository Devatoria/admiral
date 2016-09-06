package auth

import (
	"errors"
	"net/http"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/models"

	"golang.org/x/crypto/bcrypt"
)

// Authenticate returns an error if the given request basic auth is unable to authenticate the user, nil otherwise
func Authenticate(req *http.Request) error {
	username, password, ok := req.BasicAuth()
	if !ok {
		return errors.New("Unable to get request basic auth")
	}

	var user models.User
	db.Instance().Where("username = ?", username).Find(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err
}
