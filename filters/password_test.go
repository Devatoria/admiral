package filters

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	var password string

	// Nominal
	password = "P@ssw0rd#!/+=_-$€£&()!§"
	if ValidatePassword(password) != nil {
		t.Fail()
	}

	// Too short
	password = "abc"
	if ValidatePassword(password) == nil {
		t.Fail()
	}

	// Too long
	password = ""
	for i := 0; i <= 64; i++ {
		password += "a"
	}
	if ValidatePassword(password) == nil {
		t.Fail()
	}

	// Invalid chars
	password = "abcdèç'`"
	if ValidatePassword(password) == nil {
		t.Fail()
	}
}
