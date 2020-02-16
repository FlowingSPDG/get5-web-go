package config

import (
	"flag"
	"github.com/go-ini/ini"
	"strings"
	"testing"
)

// Config Configration Struct for config.ini
type Config struct {
	SteamAPIKey      string
	DefaultPage      string
	SQLHost          string
	SQLUser          string
	SQLPass          string
	SQLPort          int
	SQLDBName        string
	SQLDebugMode     bool
	HOST             string
	Cookie           string
	APIONLY          bool
	GrpcAddr         string
	EnablegRPC       bool
	ActiveMapPool    []string
	ReserveMapPool   []string
	UserMaxResources UserMaxResources
}

type UserMaxResources struct {
	Servers uint16
	Teams   uint16
	Matches uint16
}

var (
	// Cnf is configration struct
	Cnf        Config
	ConfigPath *string
)

func init() {
	ConfigPath = flag.String("cfg", "config.ini", "path to config.ini")
	testing.Init()
	flag.Parse()
	c, err := ini.Load(*ConfigPath)
	if err != nil {
		panic(err)
	}
	active := c.Section("MAPLIST").Key("Active").MustString("de_dust2,de_mirage,de_inferno,de_overpass,de_train,de_nuke,de_vertigo")
	reserve := c.Section("MAPLIST").Key("Reserve").MustString("de_cache,de_season")
	Cnf = Config{
		SteamAPIKey:    c.Section("Steam").Key("APIKey").MustString(""),
		DefaultPage:    c.Section("GET5").Key("DefaultPage").MustString(""),
		HOST:           c.Section("GET5").Key("HOST").MustString(""),
		SQLHost:        c.Section("sql").Key("host").MustString(""),
		SQLUser:        c.Section("sql").Key("user").MustString(""),
		SQLPass:        c.Section("sql").Key("pass").MustString(""),
		SQLPort:        c.Section("sql").Key("port").MustInt(3306),
		SQLDBName:      c.Section("sql").Key("database").MustString(""),
		SQLDebugMode:   c.Section("sql").Key("debug").MustBool(false),
		Cookie:         c.Section("GET5").Key("Cookie").MustString("SecretString"),
		APIONLY:        c.Section("GET5").Key("API_ONLY").MustBool(false),
		GrpcAddr:       c.Section("GET5").Key("GPRC_ADDR").MustString(":50055"),
		EnablegRPC:     c.Section("GET5").Key("ENABLE_gRPC").MustBool(false),
		ActiveMapPool:  strings.Split(strings.ToLower(strings.TrimSpace(active)), ","),
		ReserveMapPool: strings.Split(strings.ToLower(strings.TrimSpace(reserve)), ","),
		UserMaxResources: UserMaxResources{
			Servers: uint16(c.Section("USER").Key("Max_Servers").MustUint(10)),
			Teams:   uint16(c.Section("USER").Key("Max_Teams").MustUint(100)),
			Matches: uint16(c.Section("USER").Key("Max_Matches").MustUint(1000)),
		},
	}

}
