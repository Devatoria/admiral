package filters

import (
	"errors"
	"regexp"
	"strings"
)

var usernameRegex *regexp.Regexp

func init() {
	var err error
	usernameRegex, err = regexp.Compile("^[a-zA-Z0-9_-]{3,64}$")
	if err != nil {
		panic(err)
	}
}

// SanitizeUsername ensure that username fit requirements
func SanitizeUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	if !usernameRegex.MatchString(username) {
		return "", errors.New("Username does not match conditions")
	}

	return username, nil
}
