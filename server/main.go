package main

import (
	"log"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/FlowingSPDG/csgo-log-http"
	v1api "github.com/FlowingSPDG/get5-web-go/server/src/api/v1"
	v2api "github.com/FlowingSPDG/get5-web-go/server/src/api/v2"
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

func meHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*db.UserData).Name,
		"admin":    user.(*db.UserData).Admin,
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
		v1.GET("/GetMatches", v1api.GetMatches)
		v1.GET("/GetMetrics", v1api.GetMetrics)
		v1.GET("/GetSteamName", v1api.GetSteamName)
		v1.GET("/GetTeamList", v1api.GetTeamList)
		v1.GET("/GetServerList", v1api.GetServerList)
		v1.GET("/GetVersion", v1api.GetVersion)
		v1.GET("/GetMapList", v1api.GetMapList)
		v1.GET("/CheckLoggedIn", v1api.CheckLoggedIn)

		// v1 API(Legacy, authenticated with sessions.)
		match := v1.Group("/match")
		{
			match.GET("/:matchID/GetMatchInfo", v1api.GetMatchInfo)
			match.GET("/:matchID/GetPlayerStatInfo", v1api.GetPlayerStatInfo)
			match.GET("/:matchID/GetStatusString", v1api.GetStatusString)
			match.POST("/:matchID", v1api.CreateMatch) // avoid conflicts...

			// GET5 API
			match.GET("/:matchID/config", v1api.MatchConfigHandler)
			match.POST("/:matchID/finish", v1api.MatchFinishHandler)
			match.POST("/:matchID/map/:mapNumber/start", v1api.MatchMapStartHandler)
			match.POST("/:matchID/map/:mapNumber/update", v1api.MatchMapUpdateHandler)
			match.POST("/:matchID/map/:mapNumber/finish", v1api.MatchMapFinishHandler)
			match.POST("/:matchID/map/:mapNumber/player/:steamid64/update", v1api.MatchMapPlayerUpdateHandler)

			match.POST("/:matchID/cancel", v1api.MatchCancelHandler)
			match.POST("/:matchID/rcon", v1api.MatchRconHandler)
			match.POST("/:matchID/pause", v1api.MatchPauseHandler)
			match.POST("/:matchID/unpause", v1api.MatchUnpauseHandler)
			match.POST("/:matchID/adduser", v1api.MatchAddUserHandler)
			// // match.POST("/:matchID/sendconfig", v1api.MatchSendConfigHandler) // ? // I won't implement this
			match.GET("/:matchID/backup", v1api.MatchListBackupsHandler)
			match.POST("/:matchID/backup", v1api.MatchLoadBackupsHandler)

			// match.POST("/:matchID/vetoUpdate", v1api.MatchVetoUpdateHandler)
			// match.POST("/:matchID/map/:mapNumber/demo", v1api.MatchDemoUploadHandler)

			// CSGO Server log parsing
			match.POST("/:matchID/csgolog/:auth", csgologhttp.CSGOLogger(logging.MessageHandler))
		}

		team := v1.Group("/team")
		{
			team.GET("/:teamID/GetTeamInfo", v1api.GetTeamInfo)
			team.GET("/:teamID/GetRecentMatches", v1api.GetRecentMatches)
			team.GET("/:teamID/CheckUserCanEdit", v1api.CheckUserCanEdit)
			team.POST("/create", v1api.CreateTeam)
			team.PUT("/:teamID/edit", v1api.EditTeam)
			team.DELETE("/:teamID/delete", v1api.DeleteTeam)
		}

		user := v1.Group("/user")
		{
			user.GET("/:userID/GetUserInfo", v1api.GetUserInfo)
		}

		server := v1.Group("/server")
		{
			server.GET("/:serverID/GetServerInfo", v1api.GetServerInfo)
			server.POST("/create", v1api.CreateServer)
			server.PUT("/:serverID/edit", v1api.EditServer)
			server.DELETE("/:serverID/delete", v1api.DeleteServer)
		}
	}

	// V2 API(New. authenticated with JWT.)
	v2 := r.Group("/api/v2")
	{
		v2.GET("/me", v2api.AuthMidldleware.MiddlewareFunc(), meHandler)

		/*
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
		*/
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
