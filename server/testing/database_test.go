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
