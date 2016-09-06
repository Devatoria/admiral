package filters

import (
	"errors"
	"regexp"
)

var usernameRegex *regexp.Regexp

func init() {
	var err error
	usernameRegex, err = regexp.Compile("[a-zA-Z0-9_-]{3,64}")
	if err != nil {
		panic(err)
	}
}

// ValidateUsername returns an error if the given username doesn't match the regex, nil otherwise
func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return errors.New("Username does not match conditions")
	}

	return nil
}
