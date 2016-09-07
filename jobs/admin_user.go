package jobs

import (
	"errors"

	"github.com/Devatoria/admiral/db"
	"github.com/Devatoria/admiral/filters"
	"github.com/Devatoria/admiral/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdminUser(args []string) error {
	if len(args) != 3 {
		return errors.New("Wrong number of arguments: create_admin <username> <password>")
	}

	username := args[1]
	password := args[2]

	// Check that username does not already exist
	var user models.User
	db.Instance().Where("username = ?", username).Find(&user)
	if user.ID != 0 {
		return errors.New("This user already exists")
	}

	// Check username validity
	var err error
	if err = filters.ValidateUsername(username); err != nil {
		return errors.New("Username doesn't match required conditions")
	}

	// Check password validity
	if err = filters.ValidatePassword(password); err != nil {
		return errors.New("Password doesn't match required conditions")
	}

	// Hash password before storing it in database
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	user.Username = username
	user.Password = string(hash)
	user.IsAdmin = true
	db.Instance().Create(&user)

	return nil
}
