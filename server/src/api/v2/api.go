package v2

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/solovev/steam_go"

	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
)

const (
	// VERSION get5-web-go Version
	VERSION = "0.1.3"

	identityKey = "id"
	nameKey     = "name"
	adminKey    = "admin"
)

var (
	// AuthMidldleware is middleware for steam and jwt sign-in
	AuthMidldleware *jwt.GinJWTMiddleware
)

func init() {
	// the jwt middleware

	var err error
	AuthMidldleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "GET5-WEB-GO JWT",
		Key:         []byte(config.Cnf.Cookie),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// log.Printf("data : %v\n", data)
			if v, ok := data.(*db.UserData); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					nameKey:     v.Name,
					adminKey:    v.Admin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Printf("claims : %v\n", claims)
			return &db.UserData{
				ID:    int(claims[identityKey].(float64)),
				Name:  claims[nameKey].(string),
				Admin: claims[adminKey].(bool),
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
				user := &db.UserData{}
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
			c.Redirect(302, body)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*db.UserData); ok {
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

	if err != nil {
		log.Println("FAILED TO INITIALIZE GIN-JWT MIDDLEWERE.")
		panic(err)
	}
}
