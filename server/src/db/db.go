package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"

	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/go-sessions"
)

// DBdatas Struct for MySQL configration and Gorm
type DBdatas struct {
	Host string
	User string
	Pass string
	Db   string
	Port int
	Gorm *gorm.DB
}

var (
	// SteamAPIKey Steam Web API Key for accessing Steam API.
	SteamAPIKey = ""
	// DefaultPage Default page where player access root directly.
	DefaultPage string
	// SQLAccess SQL Access Object for MySQL and GORM things
	SQLAccess DBdatas
	// Sess Session
	Sess *sessions.Sessions
)

func init() {
	SQLAccess = DBdatas{
		Host: config.Cnf.SQLHost,
		User: config.Cnf.SQLUser,
		Pass: config.Cnf.SQLPass,
		Db:   config.Cnf.SQLDBName,
		Port: config.Cnf.SQLPort,
	}
	sqloption := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", SQLAccess.User, SQLAccess.Pass, SQLAccess.Host, SQLAccess.Port, SQLAccess.Db)
	//log.Println(sqloption)
	DB, err := gorm.Open("mysql", sqloption)
	if err != nil {
		panic(err)
	}
	if config.Cnf.SQLDebugMode {
		log.Println("SQL Debug mode Enabled. Transaction logs active")
	}
	DB.LogMode(config.Cnf.SQLDebugMode)
	SQLAccess.Gorm = DB
	SteamAPIKey = config.Cnf.SteamAPIKey
	DefaultPage = config.Cnf.DefaultPage
}

// GetUserData Gets UserData array via MySQL(GORM).
func (s *DBdatas) GetUserData(limit int, wherekey string, wherevalue string) ([]UserData, error) {
	UserData := []UserData{}
	s.Gorm.Limit(limit).Where(wherekey+" = ?", wherevalue).Find(&UserData)
	return UserData, nil
}
