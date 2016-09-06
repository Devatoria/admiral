package filters

import (
	"errors"
	"regexp"
)

var teamRegex *regexp.Regexp

func init() {
	var err error
	teamRegex, err = regexp.Compile("^[a-zA-Z0-9_-]{1,64}$")
	if err != nil {
		panic(err)
	}
}

func ValidateTeam(team string) error {
	if !teamRegex.MatchString(team) {
		return errors.New("Team does not match conditions")
	}

	return nil
}
