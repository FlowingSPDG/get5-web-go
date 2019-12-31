package util

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/FlowingSPDG/get5-web-go/server/src/db"

	_ "github.com/gorilla/mux"
	_ "github.com/gorilla/sessions"
	steam "github.com/kidoman/go-steam"
	"strconv"

	// a2s "github.com/rumblefrog/go-a2s"
	_ "github.com/solovev/steam_go"
	//_ "html/template"
	_ "net/http"
	_ "strconv"
	_ "time"
)

func checkIP(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		fmt.Printf("%v is not an IPv4 address\n", ip)
		return false
	}
	return true
}

// FormatMapName Formats correct map name.
func FormatMapName(mapname string) string {
	FormattedNames := make(map[string]string)
	FormattedNames["de_dust2"] = "Dust II"
	FormattedNames["de_mirage"] = "Mirage"
	FormattedNames["de_overpass"] = "Overpass"
	FormattedNames["de_vertigo"] = "Vertigo"
	FormattedNames["de_nuke"] = "NUKE"
	FormattedNames["de_train"] = "Train"
	FormattedNames["de_inferno"] = "Inferno"

	return FormattedNames["mapname"]
}

// SendRCON Sends Remote-Commands to specific IP SRCDS.
func SendRCON(host string, port int, pass string, cmd string) (string, error) {
	if !checkIP(host) {
		return "", fmt.Errorf("Specified IP is not valid")
	}
	dest := host + ":" + strconv.Itoa(port)
	o := &steam.ConnectOptions{RCONPassword: pass}
	rcon, err := steam.Connect(dest, o)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer rcon.Close()

	resp, err := rcon.Send(cmd)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return resp, nil
}

// CheckServerConnection Check server pulse by sending "status" command
func CheckServerConnection(srv db.GameServerData) bool {
	_, err := SendRCON(srv.IPString, srv.Port, srv.RconPassword, "status")
	if err != nil {
		return false
	}
	return true
}

// GET5AvailableDatas Struct for Get5 availability check.
type GET5AvailableDatas struct {
	Gamestate     int    `json:"gamestate"`
	Available     int    `json:"available"`
	PluginVersion string `json:"plugin_version"`
}

// CheckServerAvailability if server is usable for get5_web
func CheckServerAvailability(srv *db.GameServerData) (bool, string) { // available or error string
	resp, err := SendRCON(srv.IPString, srv.Port, srv.RconPassword, "get5_web_avaliable")
	if err != nil {
		return false, "Connect fails"
	}
	jsonBytes := ([]byte)(resp)
	data := new(GET5AvailableDatas)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return false, "Error reading get5_web_avaliable response"
	}
	if strings.Contains(resp, "Unknown command") {
		return false, "Either get5 or get5_apistats plugin missin"
	}
	if data.Gamestate != 0 {
		return false, "Server already has a get5 match setup"
	}
	return true, ""

}

func GetVersion() (string, error) {
	//root_dir, err := os.Executable()
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "", nil
	}
	return string(out), nil
}
