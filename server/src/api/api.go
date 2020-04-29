package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/solovev/steam_go"

	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
	"github.com/FlowingSPDG/get5-web-go/server/src/util"
)

const (
	// VERSION get5-web-go Version
	VERSION = "0.1.3"

	identityKey = "id"
	nameKey     = "name"
)

var (
	// AuthMidldleware is middleware for steam and jwt sign-in
	AuthMidldleware *jwt.GinJWTMiddleware
)

func init() {
	// the jwt middleware

	var err error
	AuthMidldleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
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
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Printf("claims : %v\n", claims)
			return &db.UserData{
				ID:   int(claims[identityKey].(float64)),
				Name: claims[nameKey].(string),
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

// GetVersion handler for /api/v1/GetVersion API.
func GetVersion(c *gin.Context) {
	log.Println("GetVersion")
	c.JSON(http.StatusOK, gin.H{
		"version": VERSION,
	})
}

// GetMapList handler for /api/v1/GetMapList API.
func GetMapList(c *gin.Context) {
	log.Println("GetMapList")
	c.JSON(http.StatusOK, gin.H{
		"active":  config.Cnf.ActiveMapPool,
		"reserve": config.Cnf.ReserveMapPool,
	})
}

// CheckLoggedIn handler for /api/v1/CheckLoggedIn API.
func CheckLoggedIn(c *gin.Context) {
	log.Println("CheckLoggedIn")
	identityKey := "id"
	claims := jwt.ExtractClaims(c)
	log.Printf("JWT_TOKEN : %v", c.Keys["JWT_TOKEN"])
	user, exist := c.Get(identityKey)
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"isLoggedIn": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isLoggedIn": true,
		// "isAdmin":    s.Get("Admin").(bool),
		// "steamid":    s.Get("SteamID").(string),
		"user":   user.(*db.UserData),
		"userid": claims[identityKey],
	})
}

// GetMatchInfo Gets match info by ID
func GetMatchInfo(c *gin.Context) {
	log.Printf("GetMatchInfo\n")
	matchid := c.Params.ByName("matchID")
	match := db.MatchData{}
	mapstats := []APIMapStatsData{}
	server := APIGameServerData{}
	team1 := APITeamData{}
	team2 := APITeamData{}
	user := APIUserData{}
	db.SQLAccess.Gorm.First(&match, matchid)
	db.SQLAccess.Gorm.Where("match_id = ?", matchid).Limit(7).Find(&mapstats)
	db.SQLAccess.Gorm.First(&server, match.ServerID)
	db.SQLAccess.Gorm.First(&team1, match.Team1ID)
	db.SQLAccess.Gorm.First(&team2, match.Team2ID)
	db.SQLAccess.Gorm.First(&user, match.UserID)
	var winner int64
	if match.Winner.Valid {
		winner_, err := match.Winner.Value()
		if err != nil {
			winner = 0
		}
		winner = winner_.(int64)
	}
	starttime := match.StartTime.Time
	endtime := match.EndTime.Time
	status, err := match.GetStatusString(true)
	if err != nil {
		status = ""
	}

	for i := 0; i < len(mapstats); i++ {
		if mapstats[i].StartTime.Valid {
			mapstats[i].StartTimeJSON = mapstats[i].StartTime.Time
		}
		if mapstats[i].EndTime.Valid {
			mapstats[i].EndTimeJSON = mapstats[i].EndTime.Time
		}
	}
	response := APIMatchData{
		ID:            match.ID,
		UserID:        match.UserID,
		Winner:        winner,
		Cancelled:     match.Cancelled,
		StartTimeJSON: starttime,
		EndTimeJSON:   endtime,
		MaxMaps:       match.MaxMaps,
		Title:         match.Title,
		SkipVeto:      match.SkipVeto,
		VetoMapPool:   strings.Split(match.VetoMapPool, " "),
		Team1Score:    match.Team1Score,
		Team2Score:    match.Team2Score,
		Team1String:   match.Team1String,
		Team2String:   match.Team2String,
		Forfeit:       match.Forfeit,
		Pending:       match.Pending(),
		Live:          match.Live(),
		Server:        server,
		MapStats:      mapstats,
		Team1:         team1,
		Team2:         team2,
		User:          user,
		Status:        status,
	}
	c.JSON(http.StatusOK, response)
}

func GetPlayerStatInfo(c *gin.Context) {
	log.Printf("GetPlayerStatInfo\n")
	matchID := c.Params.ByName("matchID")
	mapID := c.Query("mapID")
	response := []APIPlayerStatsData{}
	db.SQLAccess.Gorm.Where("match_id = ? AND map_id = ?", matchID, mapID).Limit(10).Find(&response)
	for i := 0; i < len(response); i++ { // Calculates by server-side for avoiding JavaScript's float restrcition
		response[i].Rating = response[i].GetRating()
		response[i].KDR = response[i].GetKDR()
		response[i].HSP = response[i].GetHSP()
		response[i].ADR = response[i].GetADR()
		response[i].FPR = response[i].GetFPR()
	}
	c.JSON(http.StatusOK, response)
}

func GetUserInfo(c *gin.Context) {
	log.Printf("GetUserInfo\n")
	userid := c.Params.ByName("userID")
	user := &APIUserData{}
	db.SQLAccess.Gorm.First(user, userid)

	db.SQLAccess.Gorm.Model(user).Related(&user.Teams, "Teams")
	db.SQLAccess.Gorm.Model(user).Related(&user.Servers, "Servers")

	matches := []db.MatchData{}

	db.SQLAccess.Gorm.Model(user).Related(&matches, "Matches")

	var m []*APIMatchData
	for i := 0; i < len(matches); i++ {
		mapstats := []APIMapStatsData{}
		server := APIGameServerData{}
		team1 := APITeamData{}
		team2 := APITeamData{}
		user := APIUserData{}
		db.SQLAccess.Gorm.Where("match_id = ?", matches[i].ID).Limit(7).Find(&mapstats)
		db.SQLAccess.Gorm.First(&server, matches[i].ServerID)
		db.SQLAccess.Gorm.First(&user, matches[i].UserID)
		db.SQLAccess.Gorm.First(&team1, matches[i].Team1ID)
		db.SQLAccess.Gorm.First(&team2, matches[i].Team2ID)
		var winner int64
		if matches[i].Winner.Valid {
			winner_, err := matches[i].Winner.Value()
			if err != nil {
				winner = 0
			}
			winner = winner_.(int64)
		}
		starttime := matches[i].StartTime.Time
		endtime := matches[i].EndTime.Time
		status, err := matches[i].GetStatusString(true)
		if err != nil {
			status = ""
		}

		for i := 0; i < len(mapstats); i++ {
			if mapstats[i].StartTime.Valid {
				mapstats[i].StartTimeJSON = mapstats[i].StartTime.Time
			}
			if mapstats[i].EndTime.Valid {
				mapstats[i].EndTimeJSON = mapstats[i].EndTime.Time
			}
		}
		matchdata := APIMatchData{
			ID:            matches[i].ID,
			UserID:        matches[i].UserID,
			Winner:        winner,
			Cancelled:     matches[i].Cancelled,
			StartTimeJSON: starttime,
			EndTimeJSON:   endtime,
			MaxMaps:       matches[i].MaxMaps,
			Title:         matches[i].Title,
			SkipVeto:      matches[i].SkipVeto,
			VetoMapPool:   strings.Split(matches[i].VetoMapPool, " "),
			Team1Score:    matches[i].Team1Score,
			Team2Score:    matches[i].Team2Score,
			Team1String:   matches[i].Team1String,
			Team2String:   matches[i].Team2String,
			Forfeit:       matches[i].Forfeit,
			Pending:       matches[i].Pending(),
			Live:          matches[i].Live(),
			Server:        server,
			MapStats:      mapstats,
			Team1:         team1,
			Team2:         team2,
			User:          user,
			Status:        status,
		}
		m = append(m, &matchdata)
	}
	user.Matches = m
	c.JSON(http.StatusOK, user)
}

func GetServerInfo(c *gin.Context) {
	log.Printf("GetServerInfo\n")
	serverID := c.Params.ByName("serverID")
	response := db.GameServerData{}
	db.SQLAccess.Gorm.First(&response, serverID)
	c.JSON(http.StatusOK, response)
}

func GetStatusString(c *gin.Context) {
	log.Printf("GetStatusString\n")
	matchid := c.Params.ByName("matchID")
	response := db.MatchData{}
	db.SQLAccess.Gorm.First(&response, matchid)
	status, err := response.GetStatusString(true)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, status)
}

func GetMatches(c *gin.Context) {
	log.Printf("GetMatches\n")
	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}
	userID := c.Query("userID")
	user := db.UserData{}
	response := []db.MatchData{}
	if userID != "" {
		db.SQLAccess.Gorm.Limit(20).First(&user, userID)
		db.SQLAccess.Gorm.Model(&user).Order("id DESC").Related(&response, "Matches")
	} else {
		db.SQLAccess.Gorm.Limit(20).Order("id DESC").Offset(offset).Find(&response)
	}
	c.JSON(http.StatusOK, response)
}

func GetTeamInfo(c *gin.Context) {
	log.Printf("GetTeamInfo\n")
	teamid, err := strconv.Atoi(c.Params.ByName("teamID"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("teamID should be int"))
		return
	}
	response := APITeamData{}
	rec := db.SQLAccess.Gorm.First(&response, teamid)
	if rec.RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("team not found"))
		return
	}
	var steamids = make([]string, 5)
	steamids, err = response.GetPlayers()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i := 0; i < len(steamids); i++ {
		response.SteamIDs = append(response.SteamIDs, steamids[i])
	}
	c.JSON(http.StatusOK, response)
}

func GetRecentMatches(c *gin.Context) {
	log.Printf("GetRecentMatches\n")
	teamID := c.Params.ByName("teamID")
	matches := []db.MatchData{}
	response := make([]APIMatchData, 0, len(matches))
	db.SQLAccess.Gorm.Where("team1_id = ?", teamID).Or("team2_id = ?", teamID).Limit(20).Order("id DESC").Find(&matches)
	for i := 0; i < len(matches); i++ {
		mapstats := []APIMapStatsData{}
		server := APIGameServerData{}
		team1 := APITeamData{}
		team2 := APITeamData{}
		user := APIUserData{}
		db.SQLAccess.Gorm.Where("match_id = ?", matches[i].ID).Limit(7).Find(&mapstats)
		db.SQLAccess.Gorm.First(&server, matches[i].ServerID)
		db.SQLAccess.Gorm.First(&team1, matches[i].Team1ID)
		db.SQLAccess.Gorm.First(&team2, matches[i].Team2ID)
		db.SQLAccess.Gorm.First(&user, matches[i].UserID)
		var winner int64
		if matches[i].Winner.Valid {
			winner_, err := matches[i].Winner.Value()
			if err != nil {
				winner = 0
			}
			winner = winner_.(int64)
		}
		starttime := matches[i].StartTime.Time
		endtime := matches[i].EndTime.Time
		status, err := matches[i].GetStatusString(true)
		if err != nil {
			status = ""
		}

		for i := 0; i < len(mapstats); i++ {
			if mapstats[i].StartTime.Valid {
				mapstats[i].StartTimeJSON = mapstats[i].StartTime.Time
			}
			if mapstats[i].EndTime.Valid {
				mapstats[i].EndTimeJSON = mapstats[i].EndTime.Time
			}
		}
		m := APIMatchData{
			ID:            matches[i].ID,
			UserID:        matches[i].UserID,
			Winner:        winner,
			Cancelled:     matches[i].Cancelled,
			StartTimeJSON: starttime,
			EndTimeJSON:   endtime,
			MaxMaps:       matches[i].MaxMaps,
			Title:         matches[i].Title,
			SkipVeto:      matches[i].SkipVeto,
			VetoMapPool:   strings.Split(matches[i].VetoMapPool, " "),
			Team1Score:    matches[i].Team1Score,
			Team2Score:    matches[i].Team2Score,
			Team1String:   matches[i].Team1String,
			Team2String:   matches[i].Team2String,
			Forfeit:       matches[i].Forfeit,
			Pending:       matches[i].Pending(),
			Live:          matches[i].Live(),
			Server:        server,
			MapStats:      mapstats,
			Team1:         team1,
			Team2:         team2,
			User:          user,
			Status:        status,
		}
		response = append(response, m)
	}
	c.JSON(http.StatusOK, response)
}

func CheckUserCanEdit(c *gin.Context) {
	log.Printf("CheckUserCanEdit\n")
	teamid := c.Params.ByName("teamID")
	useridstr := c.Query("userID")
	team := db.TeamData{}
	db.SQLAccess.Gorm.First(&team, teamid)
	userid, err := strconv.Atoi(useridstr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("userid should be int"))
		return
	}
	c.String(http.StatusOK, strconv.FormatBool(team.CanEdit(userid)))
}

func GetMetrics(c *gin.Context) {
	log.Printf("GetMetrics\n")
	Metrics := db.GetMetrics()
	c.JSON(http.StatusOK, Metrics)
}

// GetSteamName Get Steam Profile name by SteamWebAPI
func GetSteamName(c *gin.Context) {
	log.Printf("GetSteamName\n")
	steamid := c.Query("steamID")
	steamid64, err := strconv.Atoi(steamid)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("steamID should be int"))
		return
	}
	steamname, err := db.GetSteamName(uint64(steamid64))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, steamname)
}

// GetTeamList Returns registered team list in JSON.
func GetTeamList(c *gin.Context) {
	log.Printf("GetTeamList\n")
	Teams := []db.TeamData{}
	s := db.Sess.Start(c.Writer, c.Request)
	var IsLoggedIn bool
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		db.SQLAccess.Gorm.Where("public_team = 1").Or("user_id = ?", userid).Find(&Teams)
	} else {
		db.SQLAccess.Gorm.Where("public_team = 1").Find(&Teams)
	}
	for i := 0; i < len(Teams); i++ {
		Teams[i].GetPlayers()
	}
	c.JSON(http.StatusOK, Teams)
}

// GetServerList Returns registered public server and owned list in JSON
func GetServerList(c *gin.Context) {
	log.Printf("GetServerList\n")
	s := db.Sess.Start(c.Writer, c.Request)
	var IsLoggedIn bool
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	Servers := []APIGameServerData{}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		db.SQLAccess.Gorm.Where("public_server = true AND in_use = false").Or("user_id = ? AND in_use = false", userid).Find(&Servers)
	} else {
		db.SQLAccess.Gorm.Where("public_server = true AND in_use = false").Find(&Servers)
	}
	c.JSON(http.StatusOK, Servers)
}

// CreateTeam Registers team info to DB
func CreateTeam(c *gin.Context) {
	log.Printf("CreateTeam\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		Team := db.TeamData{}
		TeamTemp := db.TeamData{}
		err := json.NewDecoder(c.Request.Body).Decode(&TeamTemp)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("JSON Format invalid %s", err.Error()))
			return
		}
		_, err = Team.Create(userid, TeamTemp.Name, TeamTemp.Tag, TeamTemp.Flag, TeamTemp.Logo, TeamTemp.SteamIDsJSON, TeamTemp.PublicTeam)
		if err != nil {
			log.Printf("Failed to create team. ERR : %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("Forbidden"))
	}
}

// EditTeam Edits team information
func EditTeam(c *gin.Context) {
	log.Printf("EditTeam\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		teamid, err := strconv.Atoi(c.Params.ByName("teamID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("teamid should be int."))
			return
		}
		Team := db.TeamData{}
		db.SQLAccess.Gorm.First(&Team, teamid)
		if !Team.CanEdit(userid) {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("You do not have permission to edit this team"))
			return
		}
		err = json.NewDecoder(c.Request.Body).Decode(&Team)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("JSON Format Invalid %s", err.Error()))
			return
		}
		Team.ID = teamid
		Team.UserID = userid
		log.Printf("Team : %v\n", Team)

		_, err = Team.Edit()
		if err != nil {
			log.Printf("Team edit ERR : %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// DeleteTeam Deletes team
func DeleteTeam(c *gin.Context) {
	log.Printf("DeleteTeam\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		teamID, err := strconv.Atoi(c.Params.ByName("teamID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("teamID should be int."))
			return
		}
		Team := db.TeamData{ID: teamID}

		if !Team.CanDelete(userid) {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("You dont have permission to delete this team"))
			return
		}

		err = Team.Delete()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// CreateServer Register server to DB
func CreateServer(c *gin.Context) {
	log.Printf("CreateServer\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		Server := db.GameServerData{}
		ServerTemp := db.GameServerData{}
		err := json.NewDecoder(c.Request.Body).Decode(&ServerTemp)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("JSON Format invalid : %s", err.Error()))
			return
		}
		_, err = Server.Create(userid, ServerTemp.DisplayName, ServerTemp.IPString, ServerTemp.Port, ServerTemp.RconPassword, ServerTemp.PublicServer)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to create server"))
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// EditServer Edits Server information
func EditServer(c *gin.Context) {
	log.Printf("EditServer\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		serverID, err := strconv.Atoi(c.Params.ByName("serverID"))
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("serverID should be int"))
			return
		}
		Server := db.GameServerData{ID: serverID}
		if !Server.CanEdit(userid) {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("You do not have permission to edit this server"))
			return
		}
		err = json.NewDecoder(c.Request.Body).Decode(&Server)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("JSON Format invalid %s", err.Error()))
			return
		}
		Server.ID = serverID
		Server.UserID = userid

		_, err = Server.Edit()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// DeleteServer Deletes Server information
func DeleteServer(c *gin.Context) {
	log.Printf("DeleteTeam\n")
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		serverID, err := strconv.Atoi(c.Params.ByName("serverID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("serverID should be int"))
			return
		}
		Server := db.GameServerData{ID: serverID}

		if !Server.CanDelete(userid) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("You dont have permission to delete this server"))
			return
		}

		err = Server.Delete()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// CreateMatch Registers match info
func CreateMatch(c *gin.Context) {
	log.Printf("CreateMatch\n")
	if c.Params.ByName("matchID") != "create" { // rejects requests other than "/match/create".
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Unknown Request"))
		return
	}
	var IsLoggedIn = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
	}
	if IsLoggedIn {
		userid := s.Get("UserID").(int)
		var MatchTemp = db.MatchData{}
		var Match = db.MatchData{}
		err := json.NewDecoder(c.Request.Body).Decode(&MatchTemp)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("JSON Format invalid %s", err.Error()))
			return
		}
		log.Printf("CVARS : %v\n", MatchTemp.CvarsJSON)
		_, err = Match.Create(userid, MatchTemp.Team1ID, MatchTemp.Team2ID, MatchTemp.Team1String, MatchTemp.Team2String, MatchTemp.MaxMaps, MatchTemp.SkipVeto, MatchTemp.Title, MatchTemp.VetoMapPoolJSON, MatchTemp.ServerID, MatchTemp.CvarsJSON, MatchTemp.SideType, MatchTemp.IsPug)
		if err != nil {
			log.Printf("ERR : %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("Forbidden"))
	}
}

// MatchCancelHandler Handler for /api/v1/match/{matchID}/cancel API.
func MatchCancelHandler(c *gin.Context) {
	log.Println("MatchCancelHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchid should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}
		MatchUpdate := Match
		db.SQLAccess.Gorm.First(&MatchUpdate)
		MatchUpdate.Cancelled = true
		db.SQLAccess.Gorm.Model(&Match).Update(&MatchUpdate)
		db.SQLAccess.Gorm.Save(&MatchUpdate)

		Server := db.GameServerData{}
		db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		ServerUpdate := Server
		db.SQLAccess.Gorm.First(&ServerUpdate)
		ServerUpdate.InUse = false
		db.SQLAccess.Gorm.Model(&Server).Update(&ServerUpdate)
		db.SQLAccess.Gorm.Save(&ServerUpdate)

		_, err = Server.SendRcon("get5_endmatch")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
		return
	}
}

// MatchRconHandler Handler for /api/v1/match/{matchID}/rcon API.
func MatchRconHandler(c *gin.Context) {
	log.Println("MatchRconHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchid should be int"))
			return
		}
		command := c.Query("command")
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		RconRes := "No output"
		res, err := Server.SendRcon(command)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if res != "" {
			RconRes = res
		}
		c.String(http.StatusOK, RconRes)
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
		return
	}
}

// MatchPauseHandler Handler for /api/v1/match/{matchID}/pause API.
func MatchPauseHandler(c *gin.Context) {
	log.Println("MatchPauseHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchid should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		_, err = Server.SendRcon("sm_pause")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// MatchUnpauseHandler Handler for /api/v1/match/{matchID}/unpause API.
func MatchUnpauseHandler(c *gin.Context) {
	log.Println("MatchUnpauseHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchid should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		_, err = Server.SendRcon("sm_unpause")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "OK")
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// MatchAddUserHandler Handler for /api/v1/match/{matchID}/adduser API.
func MatchAddUserHandler(c *gin.Context) {
	log.Println("MatchAddUserHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		team := c.Query("team")
		if team != "team1" && team != "team2" && team != "spec" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("No team specified"))
			return
		}
		auth := c.Query("auth")
		newauth, err := util.AuthToSteamID64(auth)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Auth format invalid."))
			return
		}
		log.Printf("auth : %s", newauth)
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchID should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		res, err := Server.SendRcon(fmt.Sprintf("get5_addplayer %s %s", newauth, team))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, res)
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// BackupListJSON Struct type for /api/v1/match/{matchID}/backup
type BackupListJSON struct {
	Files []string `json:"files"`
}

// MatchListBackupsHandler Handler for /api/v1/match/{matchID}/backup API(GET).
func MatchListBackupsHandler(c *gin.Context) {
	log.Println("MatchListBackupsHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchID should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		res, err := Server.SendRcon(fmt.Sprintf("get5_listbackups %d", matchid))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resJSON := BackupListJSON{
			Files: strings.Split(strings.TrimSpace(res), "\n"),
		}
		log.Printf("BackupFiles : %v", resJSON)
		c.JSON(http.StatusOK, resJSON)
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}

// MatchLoadBackupsHandler Handler for /api/v1/match/{matchID}/backup API(POST).
func MatchLoadBackupsHandler(c *gin.Context) {
	log.Println("MatchLoadBackupsHandler")
	var IsLoggedIn = false
	var IsAdmin = false
	s := db.Sess.Start(c.Writer, c.Request)
	if s.Get("Loggedin") != nil {
		IsLoggedIn = s.Get("Loggedin").(bool)
		IsAdmin = s.Get("Admin").(bool)
	}
	if IsLoggedIn && IsAdmin {
		file := c.Query("file")
		if file == "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("No file specified"))
			return
		}
		matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchID should be int"))
			return
		}
		Match := db.MatchData{}
		rec := db.SQLAccess.Gorm.First(&Match, matchid)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find match"))
			return
		}

		Server := db.GameServerData{}
		rec = db.SQLAccess.Gorm.First(&Server, Match.ServerID)
		if rec.RecordNotFound() {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find server"))
			return
		}
		res, err := Server.SendRcon(fmt.Sprintf("get5_loadbackup %s", file))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, res)
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Forbidden"))
	}
}
