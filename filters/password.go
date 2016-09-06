package filters

import (
	"errors"
	"regexp"
)

var passwordRegex *regexp.Regexp

func init() {
	var err error
	passwordRegex, err = regexp.Compile("[a-zA-Z0-9@#&(§!)$£€=+/_-]{4,64}")
	if err != nil {
		panic(err)
	}
}

// ValidatePassword returns an error if the given password doesn't match the regex, nil otherwise
func ValidatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		return errors.New("Password does not match conditions")
	}

	return nil
}
