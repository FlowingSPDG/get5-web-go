package db

import (
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"time"

	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/go-sessions"
	"github.com/solovev/steam_go"
	"net/http"
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
	// AuthMidldleware is middleware for steam and jwt sign-in
	AuthMidldleware *jwt.GinJWTMiddleware
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

	identityKey := "id"
	// the jwt middleware
	AuthMidldleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.Cnf.Cookie),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			log.Printf("data : %v\n", data)
			if v, ok := data.(*UserData); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &UserData{
				ID: int(claims[identityKey].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			opID := steam_go.NewOpenId(c.Request)
			switch opID.Mode() {
			case "":
				c.Redirect(302, opID.AuthUrl())
				return nil, jwt.ErrFailedAuthentication
			case "cancel":
				return nil, jwt.ErrFailedAuthentication
			default:
				user := &UserData{}
				steamid, err := opID.ValidateAndGetId()
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, err)
					return nil, jwt.ErrFailedAuthentication
				}
				user.SteamID = steamid
				user, _, err = user.GetOrCreate()
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, err)
					return nil, jwt.ErrFailedAuthentication
				}
				return user, nil
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			body := fmt.Sprintf("/#/auth?auth=%s&expire=%s", token, expire)
			// log.Printf("body : %v\n", body)
			c.Redirect(302, body)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*UserData); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

}

// GetUserData Gets UserData array via MySQL(GORM).
func (s *DBdatas) GetUserData(limit int, wherekey string, wherevalue string) ([]UserData, error) {
	UserData := []UserData{}
	s.Gorm.Limit(limit).Where(wherekey+" = ?", wherevalue).Find(&UserData)
	return UserData, nil
}
