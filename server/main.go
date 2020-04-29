package main

import (
	"log"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/FlowingSPDG/csgo-log-http"
	"github.com/FlowingSPDG/get5-web-go/server/src/api"
	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
	"github.com/FlowingSPDG/get5-web-go/server/src/grpc"
	"github.com/FlowingSPDG/get5-web-go/server/src/logging"
)

const (
	identityKey = "id"
)

var (
	// StaticDir Directly where serves static files
	StaticDir = "./static"
	// SQLAccess SQL Access Object for MySQL and GORM things
	SQLAccess db.DBdatas
)

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*db.UserData).Name,
	})
}

func main() {

	defer SQLAccess.Gorm.Close()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//s := r.Host(HOST).Subrouter() // in-case if we need vhost thing

	// V1 API.
	v1 := r.Group("/api/v1")
	{
		v1.GET("/GetMatches", api.GetMatches)
		v1.GET("/GetMetrics", api.GetMetrics)
		v1.GET("/GetSteamName", api.GetSteamName)
		v1.GET("/GetTeamList", api.GetTeamList)
		v1.GET("/GetServerList", api.GetServerList)
		v1.GET("/GetVersion", api.GetVersion)
		v1.GET("/GetMapList", api.GetMapList)
		v1.GET("/CheckLoggedIn", api.CheckLoggedIn)

		// v1 API(Legacy, authenticated with sessions.)
		match := v1.Group("/match")
		{
			match.GET("/:matchID/GetMatchInfo", api.GetMatchInfo)
			match.GET("/:matchID/GetPlayerStatInfo", api.GetPlayerStatInfo)
			match.GET("/:matchID/GetStatusString", api.GetStatusString)
			match.POST("/:matchID", api.CreateMatch) // avoid conflicts...

			// GET5 API
			match.GET("/:matchID/config", api.MatchConfigHandler)
			match.POST("/:matchID/finish", api.MatchFinishHandler)
			match.POST("/:matchID/map/:mapNumber/start", api.MatchMapStartHandler)
			match.POST("/:matchID/map/:mapNumber/update", api.MatchMapUpdateHandler)
			match.POST("/:matchID/map/:mapNumber/finish", api.MatchMapFinishHandler)
			match.POST("/:matchID/map/:mapNumber/player/:steamid64/update", api.MatchMapPlayerUpdateHandler)

			match.POST("/:matchID/cancel", api.MatchCancelHandler)
			match.POST("/:matchID/rcon", api.MatchRconHandler)
			match.POST("/:matchID/pause", api.MatchPauseHandler)
			match.POST("/:matchID/unpause", api.MatchUnpauseHandler)
			match.POST("/:matchID/adduser", api.MatchAddUserHandler)
			// // match.POST("/:matchID/sendconfig", api.MatchSendConfigHandler) // ? // I won't implement this
			match.GET("/:matchID/backup", api.MatchListBackupsHandler)
			match.POST("/:matchID/backup", api.MatchLoadBackupsHandler)

			// match.POST("/:matchID/vetoUpdate", api.MatchVetoUpdateHandler)
			// match.POST("/:matchID/map/:mapNumber/demo", api.MatchDemoUploadHandler)

			// CSGO Server log parsing
			match.POST("/:matchID/csgolog/:auth", csgologhttp.CSGOLogger(logging.MessageHandler))
		}

		team := v1.Group("/team")
		{
			team.GET("/:teamID/GetTeamInfo", api.GetTeamInfo)
			team.GET("/:teamID/GetRecentMatches", api.GetRecentMatches)
			team.GET("/:teamID/CheckUserCanEdit", api.CheckUserCanEdit)
			team.POST("/create", api.CreateTeam)
			team.PUT("/:teamID/edit", api.EditTeam)
			team.DELETE("/:teamID/delete", api.DeleteTeam)
		}

		user := v1.Group("/user")
		{
			user.GET("/:userID/GetUserInfo", api.GetUserInfo)
		}

		server := v1.Group("/server")
		{
			server.GET("/:serverID/GetServerInfo", api.GetServerInfo)
			server.POST("/create", api.CreateServer)
			server.PUT("/:serverID/edit", api.EditServer)
			server.DELETE("/:serverID/delete", api.DeleteServer)
		}
	}

	// V2 API(New. authenticated with JWT.)
	v2 := r.Group("/api/v2")
	{
		v2.GET("/hello", api.AuthMidldleware.MiddlewareFunc(), helloHandler)

		// Refresh time can be longer than token timeout
		// v2.GET("/refresh_token", api.AuthMidldleware.RefreshHandler)
		v2.GET("/login", api.AuthMidldleware.LoginHandler)
		v2.GET("/logout", api.AuthMidldleware.LogoutHandler)

		v2.GET("/GetMatches", api.GetMatches)
		v2.GET("/GetMetrics", api.GetMetrics)
		v2.GET("/GetSteamName", api.GetSteamName)
		v2.GET("/GetTeamList", api.GetTeamList)
		v2.GET("/GetServerList", api.GetServerList)
		v2.GET("/GetVersion", api.GetVersion)
		v2.GET("/GetMapList", api.GetMapList)
		v2.GET("/CheckLoggedIn", api.CheckLoggedIn)

		// Match API for front(Vue)
		match := v2.Group("/match")
		{
			match.POST("/:matchID", api.AuthMidldleware.MiddlewareFunc(), api.CreateMatch) // avoid conflicts...
			match.GET("/:matchID/GetMatchInfo", api.GetMatchInfo)
			match.GET("/:matchID/GetPlayerStatInfo", api.GetPlayerStatInfo)
			match.GET("/:matchID/GetStatusString", api.GetStatusString)

			// GET5 API(SRCDS -> API).
			match.GET("/:matchID/srcds/config", api.MatchConfigHandler)
			match.POST("/:matchID/srcds/finish", api.MatchFinishHandler)
			match.POST("/:matchID/srcds/map/:mapNumber/start", api.MatchMapStartHandler)
			match.POST("/:matchID/srcds/map/:mapNumber/update", api.MatchMapUpdateHandler)
			match.POST("/:matchID/srcds/map/:mapNumber/finish", api.MatchMapFinishHandler)
			match.POST("/:matchID/srcds/map/:mapNumber/player/:steamid64/update", api.MatchMapPlayerUpdateHandler)

			// GET5 API(API -> SRCDS).
			match.POST("/:matchID/cancel", api.AuthMidldleware.MiddlewareFunc(), api.MatchCancelHandler)
			match.POST("/:matchID/rcon", api.AuthMidldleware.MiddlewareFunc(), api.MatchRconHandler)
			match.POST("/:matchID/pause", api.AuthMidldleware.MiddlewareFunc(), api.MatchPauseHandler)
			match.POST("/:matchID/unpause", api.AuthMidldleware.MiddlewareFunc(), api.MatchUnpauseHandler)
			match.POST("/:matchID/adduser", api.AuthMidldleware.MiddlewareFunc(), api.MatchAddUserHandler)
			// // match.POST("/:matchID/sendconfig", api.MatchSendConfigHandler) // ? // I won't implement this
			match.GET("/:matchID/backup", api.AuthMidldleware.MiddlewareFunc(), api.MatchListBackupsHandler)
			match.POST("/:matchID/backup", api.AuthMidldleware.MiddlewareFunc(), api.MatchLoadBackupsHandler)

			// match.POST("/:matchID/vetoUpdate", api.MatchVetoUpdateHandler)
			// match.POST("/:matchID/map/:mapNumber/demo", api.MatchDemoUploadHandler)

			// CSGO Server log parsing
			match.POST("/:matchID/csgolog/:auth", csgologhttp.CSGOLogger(logging.MessageHandler))
		}

		team := v2.Group("/team")
		{
			team.GET("/:teamID/GetTeamInfo", api.GetTeamInfo)
			team.GET("/:teamID/GetRecentMatches", api.GetRecentMatches)
			team.GET("/:teamID/CheckUserCanEdit", api.AuthMidldleware.MiddlewareFunc(), api.CheckUserCanEdit)
			team.POST("/create", api.AuthMidldleware.MiddlewareFunc(), api.CreateTeam)
			team.PUT("/:teamID/edit", api.AuthMidldleware.MiddlewareFunc(), api.EditTeam)
			team.DELETE("/:teamID/delete", api.AuthMidldleware.MiddlewareFunc(), api.DeleteTeam)
		}

		user := v2.Group("/user")
		{
			user.GET("/:userID/GetUserInfo", api.GetUserInfo)
		}

		server := v2.Group("/server")
		{
			server.GET("/:serverID/GetServerInfo", api.GetServerInfo)
			server.POST("/create", api.AuthMidldleware.MiddlewareFunc(), api.CreateServer)
			server.PUT("/:serverID/edit", api.AuthMidldleware.MiddlewareFunc(), api.EditServer)
			server.DELETE("/:serverID/delete", api.AuthMidldleware.MiddlewareFunc(), api.DeleteServer)
		}
	}

	if !config.Cnf.APIONLY {
		entrypoint := "./static/index.html"
		r.GET("/", func(c *gin.Context) { c.File(entrypoint) })
		r.Use(static.Serve("/css", static.LocalFile("./static/css", false)))
		r.Use(static.Serve("/js", static.LocalFile("./static/js", false)))
		r.Use(static.Serve("/img", static.LocalFile("./static/img", false)))
		r.Use(static.Serve("/fonts", static.LocalFile("./static/fonts", false)))

	} else {
		log.Println("API ONLY MODE")
	}

	if config.Cnf.EnablegRPC {
		log.Println("EnableGRPC option enabled. Starting gRPC server...")
		go func() {
			err := get5grpc.StartGrpc(config.Cnf.GrpcAddr)
			if err != nil {
				panic(err)
			}
		}()
	}

	log.Panicf("Failed to listen port %s : %v\n", config.Cnf.HOST, r.Run(config.Cnf.HOST))

}
