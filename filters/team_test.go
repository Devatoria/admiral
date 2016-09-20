package filters

import (
	"testing"
)

func ValidateTeamTest(t *testing.T) {
	var team string

	// Nominal
	team = "devatoria"
	if ValidateTeam(team) != nil {
		t.Fail()
	}

	// Too short
	team = ""
	if ValidateTeam == nil {
		t.Fail()
	}

	// Too long
	team = ""
	for i := 0; i <= 64; i++ {
		team += "a"
	}
	if ValidateTeam(team) == nil {
		t.Fail()
	}

	// Incorrect chars
	team = "devatoria/subteam"
	if ValidateTeam(team) == nil {
		t.Fail()
	}
}
