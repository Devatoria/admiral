package filters

import (
	"testing"
)

func TestSanitizeUsername(t *testing.T) {
	var username, expected, sanitized string
	var err error

	expected = "_devatoria-1234_"

	// Nominal case
	username = "_Devatoria-1234_"
	sanitized, err = SanitizeUsername(username)
	if err != nil || sanitized != expected {
		t.Fail()
	}

	// Correct but to trim case
	username = " _Devatoria-1234_ "
	sanitized, err = SanitizeUsername(username)
	if err != nil || sanitized != expected {
		t.Fail()
	}

	// Incorrect
	username = "Devatoria/1234"
	_, err = SanitizeUsername(username)
	if err == nil {
		t.Fail()
	}
}
