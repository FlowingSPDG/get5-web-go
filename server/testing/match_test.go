package testing

import (
	_ "github.com/FlowingSPDG/get5-web-go/server/src/api"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
	_ "github.com/FlowingSPDG/get5-web-go/server/src/grpc"

	"testing"
)

func TestCreateMatch(t *testing.T) {
	// go test -v -run TestCreateMatch -args -cfg ../config.ini
	match := &db.MatchData{}
	cvars := make(map[string]string)
	cvars["hostname"] = "GET5"
	match, err := match.Create(0, 0, 0, "TEAM1STR", "TEAM2STR", 1, false, "MATCH_TITLE", []string{"de_dust2", "de_mirage", "de_inferno"}, 0, cvars, "standard", false)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Match : %v\n", match)
}
