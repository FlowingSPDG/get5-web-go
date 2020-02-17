package testing

import (
	_ "github.com/FlowingSPDG/get5-web-go/server/src/api"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
	_ "github.com/FlowingSPDG/get5-web-go/server/src/grpc"

	"testing"
)

func TestRegisterUser(t *testing.T) {
	// go test -v -run TestRegisterUser -args -cfg ../config.ini
	user := &db.UserData{SteamID: "76561198072054549"} // FlowingSPDG
	user, exist, err := user.GetOrCreate()
	if err != nil {
		t.Fatal(err.Error())
	}
	if exist {
		t.Log("USER EXIST!")
	}
	t.Logf("REGISTERED USER : %v\n", user)
}

func TestGetTeamsByUser(t *testing.T) {
	// go test -v -run TestGetTeamsByUser -args -cfg ../config.ini
	user := &db.UserData{SteamID: "76561198072054549"} // FlowingSPDG
	user, _, err := user.GetOrCreate()
	if err != nil {
		t.Fatal(err.Error())
	}
	teams := user.GetTeams(20)
	t.Logf("Found %d teams\n", len(teams))
	for k, v := range teams {
		t.Logf("Team %d : %v\n", k, v)
	}
}

func TestGetMatchesByUser(t *testing.T) {
	// go test -v -run TestGetMatchesByUser -args -cfg ../config.ini
	user := &db.UserData{SteamID: "76561198072054549"} // FlowingSPDG
	user, _, err := user.GetOrCreate()
	if err != nil {
		t.Fatal(err.Error())
	}
	matches := user.GetRecentMatches(50)
	t.Logf("Found %d matches\n", len(matches))
	for k, v := range matches {
		t.Logf("Matche %d : %v\n", k, v)
	}
}

func TestCreateTeamByUser(t *testing.T) {
	// go test -v -run TestCreateTeamByUser -args -cfg ../config.ini
	user := &db.UserData{SteamID: "76561198072054549"} // FlowingSPDG
	user, _, err := user.GetOrCreate()
	if err != nil {
		t.Fatal(err.Error())
	}
	Team := &db.TeamData{}
	t.Log("Creating team...")
	Team, err = Team.Create(user.ID, "TEST_TEAM", "TEST_TAG", "jp", "", []string{"76561198072054549"}, false)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Created team. Team : %v\n", Team)
	err = Team.Delete()
	if err != nil {
		t.Fatal(err.Error())
	}
}
