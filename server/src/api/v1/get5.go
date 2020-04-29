package v1

import (
	"fmt"
	"github.com/FlowingSPDG/get5-web-go/server/src/db"
	pb "github.com/FlowingSPDG/get5-web-go/server/src/grpc/proto"
	pbservices "github.com/FlowingSPDG/get5-web-go/server/src/grpc/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// MatchConfigHandler Handler for /api/v1/match/{matchID}/config API.
func MatchConfigHandler(c *gin.Context) {
	log.Printf("MatchConfigHandler\n")
	matchid := c.Params.ByName("matchID")
	match := db.MatchData{}
	db.SQLAccess.Gorm.First(&match, matchid)
	res, err := match.BuildMatchDict()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// MatchAPICheck Checks if API is available or not
func MatchAPICheck(m *db.MatchData, r *http.Request) error {
	if m.APIKey != r.FormValue("key") {
		return fmt.Errorf("Wrong API Key")
	}
	if m.Finalized() {
		return fmt.Errorf("Match already finalized")
	}
	return nil
}

// MatchFinishHandler Handler for /api/v1/match/{matchID}/finish API.
func MatchFinishHandler(c *gin.Context) {
	log.Println("MatchFinishHandler")
	matchid := c.Params.ByName("matchID")
	winner := c.PostForm("winner")
	forfeit := c.PostForm("forfeit")
	log.Printf("matchid : %s\n", matchid)
	log.Printf("winner : %s\n", winner)
	log.Printf("forfeit : %s\n", forfeit)
	var Match = db.MatchData{}
	db.SQLAccess.Gorm.First(&Match, matchid)
	log.Printf("Requested API key : %s\n", c.PostForm("key"))
	log.Printf("Real API key : %s\n", Match.APIKey)
	var MatchUpdate = Match
	db.SQLAccess.Gorm.First(&MatchUpdate)
	err := MatchAPICheck(&Match, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if winner == "team1" {
		MatchUpdate.Winner.Scan(MatchUpdate.Team1ID)
	} else if winner == "team2" {
		MatchUpdate.Winner.Scan(MatchUpdate.Team2ID)
	} else {
		MatchUpdate.Winner.Scan(nil)
	}
	if forfeit == "1" {
		MatchUpdate.Forfeit = true
		if winner == "team1" {
			MatchUpdate.Team1Score = 1
			MatchUpdate.Team2Score = 0
		} else if winner == "team2" {
			MatchUpdate.Team1Score = 0
			MatchUpdate.Team2Score = 1
		} else { // added,this is not exist in original get5-web
			MatchUpdate.Cancelled = true
		}
	}
	MatchUpdate.EndTime.Scan(time.Now())
	server, err := MatchUpdate.GetServer()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Removes ALL HTTP log destinations
	server.SendRcon("logaddress_delall_http")
	serverUpdate := server
	db.SQLAccess.Gorm.First(&serverUpdate)
	serverUpdate.InUse = false
	db.SQLAccess.Gorm.Model(&server).Update(&serverUpdate)
	db.SQLAccess.Gorm.Save(&serverUpdate)
	db.SQLAccess.Gorm.Model(&Match).Update(&MatchUpdate)
	db.SQLAccess.Gorm.Save(&MatchUpdate)
	log.Printf("Finished match %v, winner = %v\n", MatchUpdate.ID, winner)

	forfeitbool, err := strconv.ParseBool(forfeit)
	if err != nil {
		log.Println("Failed to parse forfeit value,interuppting streaming")
		return
	}
	matchidInt, err := strconv.Atoi(matchid)
	if err != nil {
		log.Println("match id is not int,interuppting streaming")
		return
	}

	go pbservices.MatchesStream.Write(int32(matchidInt), &pb.MatchEventReply{
		Event: &pb.MatchEventReply_Matchfinish{
			Matchfinish: &pb.MatchEventMatchFinish{
				Matchid: int32(matchidInt),
				Winner:  winner,
				Forfeit: forfeitbool,
			},
		},
	}, true)
}

// MatchMapStartHandler Handler for /api/v1/match/{matchID}/map/{mapNumber}/start  API.
func MatchMapStartHandler(c *gin.Context) {
	log.Println("MatchMapStartHandler")
	matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchid should be int"))
		return
	}
	mapnumber, err := strconv.Atoi(c.Params.ByName("mapNumber"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("mapnumber should be int"))
		return
	}
	mapname := c.PostForm("mapname")
	log.Printf("mapname : %s\n", mapname)
	var m = db.MatchData{}
	db.SQLAccess.Gorm.First(&m, matchid)
	err = MatchAPICheck(&m, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	mUpdate := m
	db.SQLAccess.Gorm.First(&mUpdate)
	MapStats := &db.MapStatsData{}
	MapStats, err = MapStats.GetOrCreate(matchid, mapnumber, mapname)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	mUpdate.StartTime.Scan(time.Now())
	db.SQLAccess.Gorm.Model(&m).Update(&mUpdate)
	db.SQLAccess.Gorm.Save(&mUpdate)

	go pbservices.MatchesStream.Write(int32(matchid), &pb.MatchEventReply{
		Event: &pb.MatchEventReply_Mapstart{
			Mapstart: &pb.MatchEventMapStart{
				Matchid:   int32(matchid),
				Mapnumber: int32(mapnumber),
				Mapname:   mapname,
			},
		},
	}, false)
}

// MatchMapUpdateHandler Handler for /api/v1/match/{matchID}/map/{mapNumber}/update API.
func MatchMapUpdateHandler(c *gin.Context) {
	log.Println("MatchMapUpdateHandler")
	matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchID should be int"))
		return
	}
	mapnumber, err := strconv.Atoi(c.Params.ByName("mapNumber"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("mapNumber should be int"))
		return
	}
	team1score, err := strconv.Atoi(c.PostForm("team1score"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("team1score should be int"))
		return
	}
	team2score, err := strconv.Atoi(c.PostForm("team2score"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("team2score should be int"))
		return
	}
	/*
		log.Printf("matchid : %d\n", matchid)
		log.Printf("mapnumber : %d\n", mapnumber)
		log.Printf("team1score : %d\n", team1score)
		log.Printf("team2score : %d\n", team2score)
	*/

	var m = db.MatchData{}
	db.SQLAccess.Gorm.First(&m, matchid)
	err = MatchAPICheck(&m, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	MapStats := db.MapStatsData{}
	MapStatsRecord := db.SQLAccess.Gorm.Where("match_id = ? AND map_number = ?", matchid, mapnumber).First(&MapStats)
	MapStatsUpdate := MapStats
	db.SQLAccess.Gorm.First(&MapStatsUpdate)
	if !MapStatsRecord.RecordNotFound() {
		if team1score != -1 && team2score != -1 {
			MapStatsUpdate.Team1Score = team1score
			MapStatsUpdate.Team2Score = team2score
			db.SQLAccess.Gorm.Model(&MapStats).Update(&MapStatsUpdate)
			db.SQLAccess.Gorm.Save(&MapStatsUpdate)
		}
	} else {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to find map stats object"))
		return
	}

	go pbservices.MatchesStream.Write(int32(matchid), &pb.MatchEventReply{
		Event: &pb.MatchEventReply_Mapupdate{
			Mapupdate: &pb.MatchEventMapUpdate{
				Matchid:    int32(matchid),
				Mapnumber:  int32(mapnumber),
				Team1Score: int32(team1score),
				Team2Score: int32(team2score),
			},
		},
	}, false)
}

// MatchMapFinishHandler Handler for /api/v1/match/{matchID}/map/{mapNumber}/finish API.
func MatchMapFinishHandler(c *gin.Context) {
	log.Println("MatchMapFinishHandler")
	matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("matchID should be int"))
		return
	}
	mapnumber, err := strconv.Atoi(c.Params.ByName("mapNumber"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("mapNumber should be int"))
		return
	}
	winner := c.PostForm("winner")
	log.Printf("matchid : %d\n", matchid)
	log.Printf("mapnumber : %d\n", mapnumber)
	log.Printf("winner : %s\n", winner)

	m := db.MatchData{}
	db.SQLAccess.Gorm.First(&m, matchid)
	mUpdate := m
	db.SQLAccess.Gorm.First(&mUpdate)
	err = MatchAPICheck(&m, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	MapStats := db.MapStatsData{}
	MapStatsRecord := db.SQLAccess.Gorm.Where("match_id = ? AND map_number = ?", matchid, mapnumber).First(&MapStats)
	MapStatsUpdate := MapStats
	db.SQLAccess.Gorm.First(&MapStatsUpdate)
	if !MapStatsRecord.RecordNotFound() {
		MapStatsUpdate.EndTime.Scan(time.Now())
		if winner == "team1" {
			MapStatsUpdate.Winner.Scan(m.Team1ID)
			mUpdate.Team1Score++
		} else if winner == "team2" {
			MapStatsUpdate.Winner.Scan(m.Team2ID)
			mUpdate.Team2Score++
		} else {
			MapStatsUpdate.Winner.Scan(nil)
		}
		db.SQLAccess.Gorm.Model(&MapStats).Update(&MapStatsUpdate)
		db.SQLAccess.Gorm.Save(&MapStatsUpdate)
		db.SQLAccess.Gorm.Model(&m).Update(&mUpdate)
		db.SQLAccess.Gorm.Save(&mUpdate)
	} else {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to find map stats object"))
		return
	}

	go pbservices.MatchesStream.Write(int32(matchid), &pb.MatchEventReply{
		Event: &pb.MatchEventReply_Mapfinish{
			Mapfinish: &pb.MatchEventMapFinish{
				Matchid:   int32(matchid),
				Mapnumber: int32(mapnumber),
				Winner:    winner,
			},
		},
	}, false)
}

// MatchMapPlayerUpdateHandler Handler for /api/v1/match/{matchID}/map/{mapNumber}/player/{steamid64}/update API.
func MatchMapPlayerUpdateHandler(c *gin.Context) {
	log.Println("MatchMapPlayerUpdateHandler")
	matchid, err := strconv.Atoi(c.Params.ByName("matchID"))
	if err != nil {
		log.Printf("ERR : %v\n", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	mapnumber, err := strconv.Atoi(c.Params.ByName("mapNumber"))
	if err != nil {
		log.Println("mapNumber should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	steamid64 := c.Params.ByName("steamid64")

	FormName := c.PostForm("name")
	FormTeam := c.PostForm("team")
	FormKills, err := strconv.Atoi(c.PostForm("kills"))
	if err != nil {
		log.Println("kills should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormAssists, err := strconv.Atoi(c.PostForm("assists"))
	if err != nil {
		log.Println("assists should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormDeaths, err := strconv.Atoi(c.PostForm("deaths"))
	if err != nil {
		log.Println("deaths should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormFlashbangAssists, err := strconv.Atoi(c.PostForm("flashbang_assists"))
	if err != nil {
		log.Println("flashbang_assists should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormTeamKills, err := strconv.Atoi(c.PostForm("teamkills"))
	if err != nil {
		log.Println("teamkills should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormSuicides, err := strconv.Atoi(c.PostForm("suicides"))
	if err != nil {
		log.Println("suicides should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormDamage, err := strconv.Atoi(c.PostForm("damage"))
	if err != nil {
		log.Println("damage should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormHeadShotKills, err := strconv.Atoi(c.PostForm("headshot_kills"))
	if err != nil {
		log.Println("headshot_kills should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormRoundsPlayed, err := strconv.Atoi(c.PostForm("roundsplayed"))
	if err != nil {
		log.Println("roundsplayed should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormBombPlants, err := strconv.Atoi(c.PostForm("bomb_plants"))
	if err != nil {
		log.Println("bomb_plants should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormBombDefuses, err := strconv.Atoi(c.PostForm("bomb_defuses"))
	if err != nil {
		log.Println("bomb_defuses should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Form1KillRounds, err := strconv.Atoi(c.PostForm("1kill_rounds"))
	if err != nil {
		log.Println("1kill_rounds should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Form2KillRounds, err := strconv.Atoi(c.PostForm("2kill_rounds"))
	if err != nil {
		log.Println("2kill_rounds should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Form3KillRounds, err := strconv.Atoi(c.PostForm("3kill_rounds"))
	if err != nil {
		log.Println("3kill_rounds should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Form4KillRounds, err := strconv.Atoi(c.PostForm("4kill_rounds"))
	if err != nil {
		log.Println("4kill_rounds should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Form5KillRounds, err := strconv.Atoi(c.PostForm("5kill_rounds"))
	if err != nil {
		log.Println("5kill_rounds should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormV1, err := strconv.Atoi(c.PostForm("v1"))
	if err != nil {
		log.Println("v1 should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormV2, err := strconv.Atoi(c.PostForm("v2"))
	if err != nil {
		log.Println("v2 should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormV3, err := strconv.Atoi(c.PostForm("v3"))
	if err != nil {
		log.Println("v3 should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormV4, err := strconv.Atoi(c.PostForm("v4"))
	if err != nil {
		log.Println("v4 should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormV5, err := strconv.Atoi(c.PostForm("v5"))
	if err != nil {
		log.Println("v5 should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormFirstKillT, err := strconv.Atoi(c.PostForm("firstkill_t"))
	if err != nil {
		log.Println("firstkill_t should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormFirstKillCT, err := strconv.Atoi(c.PostForm("firstkill_ct"))
	if err != nil {
		log.Println("firstkill_ct should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormFirstDeathT, err := strconv.Atoi(c.PostForm("firstdeath_t"))
	if err != nil {
		log.Println("firstdeath_t should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	FormFirstDeathCT, err := strconv.Atoi(c.PostForm("firstdeath_ct"))
	if err != nil {
		log.Println("firstdeath_ct should be int")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	/*
		Form, err := strconv.Atoi(c.PostForm("tradekill")) // https://github.com/FlowingSPDG/get5-webapi/blob/e41ac0ab3c698ed67dbadcd667e55feef403e074/scripting/get5_apistats.sp#L429
		if err != nil {
			log.Println("tradekill should be int")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	*/
	/*
		log.Printf("matchid : %d\n", matchid)
		log.Printf("mapnumber : %d\n", mapnumber)
		log.Printf("key : %s\n", c.PostForm("key"))
		log.Printf("name : %s\n", FormName)
		log.Printf("team : %s\n", FormTeam)
		log.Printf("kills : %d\n", FormKills)
		log.Printf("assists : %d\n", FormAssists)
		log.Printf("deaths : %d\n", FormDeaths)
		log.Printf("flashbang_assists : %d\n", FormFlashbangAssists)
		log.Printf("teamkills : %d\n", FormTeamKills)
		log.Printf("suicides : %d\n", FormSuicides) // not working?
		log.Printf("damage : %d\n", FormDamage)
		log.Printf("headshot_kills : %d\n", FormHeadShotKills)
		log.Printf("roundsplayed : %d\n", FormRoundsPlayed)
		log.Printf("bomb_plants : %d\n", FormBombPlants)
		log.Printf("bomb_defuses : %d\n", FormBombDefuses)
		log.Printf("1kill_rounds : %d\n", Form1KillRounds)
		log.Printf("2kill_rounds : %d\n", Form2KillRounds)
		log.Printf("3kill_rounds : %d\n", Form3KillRounds)
		log.Printf("4kill_rounds : %d\n", Form4KillRounds)
		log.Printf("5kill_rounds : %d\n", Form5KillRounds)
		log.Printf("v1 : %d\n", FormV1)
		log.Printf("v2 : %d\n", FormV2)
		log.Printf("v3 : %d\n", FormV3)
		log.Printf("v4 : %d\n", FormV4)
		log.Printf("v5 : %d\n", FormV5)
		log.Printf("firstkill_t : %d\n", FormFirstKillT)
		log.Printf("firstkill_ct : %d\n", FormFirstKillCT)
		log.Printf("firstdeath_t : %d\n", FormFirstDeathT)
		log.Printf("firstdeath_ct : %d\n", FormFirstDeathCT)
	*/

	var m = db.MatchData{}
	db.SQLAccess.Gorm.First(&m, matchid)
	if m.APIKey != c.PostForm("key") {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Wrong API Key"))
		return
	}

	MapStats := db.MapStatsData{}
	MapStatsRecord := db.SQLAccess.Gorm.Where("match_id = ? AND map_number = ? ", matchid, mapnumber).First(&MapStats)
	if !MapStatsRecord.RecordNotFound() {
		p := &db.PlayerStatsData{}
		p, err := p.GetOrCreate(matchid, mapnumber, steamid64)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to solve player stats object"))
			return
		}
		pUpdate := p
		db.SQLAccess.Gorm.First(&pUpdate)

		if FormTeam == "team1" {
			p.TeamID = m.Team1ID
		} else if FormTeam == "team2" {
			p.TeamID = m.Team2ID
		}
		pUpdate.Name = FormName
		pUpdate.Kills = FormKills
		pUpdate.Assists = FormAssists
		pUpdate.Deaths = FormDeaths
		pUpdate.FlashbangAssists = FormFlashbangAssists
		pUpdate.Teamkills = FormTeamKills
		pUpdate.Suicides = FormSuicides
		pUpdate.Damage = FormDamage
		pUpdate.HeadshotKills = FormHeadShotKills
		pUpdate.Roundsplayed = FormRoundsPlayed
		pUpdate.BombPlants = FormBombPlants
		pUpdate.BombDefuses = FormBombDefuses
		pUpdate.K1 = Form1KillRounds
		pUpdate.K2 = Form2KillRounds
		pUpdate.K3 = Form3KillRounds
		pUpdate.K4 = Form4KillRounds
		pUpdate.K5 = Form5KillRounds
		pUpdate.V1 = FormV1
		pUpdate.V2 = FormV2
		pUpdate.V3 = FormV3
		pUpdate.V4 = FormV4
		pUpdate.V5 = FormV5
		pUpdate.FirstkillT = FormFirstKillT
		pUpdate.FirstkillCT = FormFirstKillCT
		pUpdate.FirstdeathT = FormFirstDeathT
		pUpdate.FirstdeathCT = FormFirstDeathCT
		db.SQLAccess.Gorm.Model(&p).Update(&pUpdate)
		db.SQLAccess.Gorm.Save(&pUpdate)
	} else {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Failed to find map stats object"))
		return
	}
	go pbservices.MatchesStream.Write(int32(matchid), &pb.MatchEventReply{
		Event: &pb.MatchEventReply_Mapplayerupdate{
			Mapplayerupdate: &pb.MatchEventMapPlayerUpdate{
				Matchid:          int32(matchid),
				Mapnumber:        int32(mapnumber),
				Steamid:          steamid64,
				Name:             FormName,
				Team:             FormTeam,
				Kills:            int32(FormKills),
				Assists:          int32(FormAssists),
				Deaths:           int32(FormDeaths),
				FlashbangAssists: int32(FormFlashbangAssists),
				Teamkills:        int32(FormTeamKills),
				Suicides:         int32(FormSuicides),
				Damage:           int32(FormDamage),
				HeadshotKills:    int32(FormHeadShotKills),
				Roundsplayed:     int32(FormRoundsPlayed),
				BombPlants:       int32(FormBombPlants),
				BombDefuses:      int32(FormBombDefuses),
				OneKillRounds:    int32(Form1KillRounds),
				TwoKillRounds:    int32(Form2KillRounds),
				ThreeKillRounds:  int32(Form3KillRounds),
				FourKillRounds:   int32(Form4KillRounds),
				FiveKillRounds:   int32(Form5KillRounds),
				V1:               int32(FormV1),
				V2:               int32(FormV2),
				V3:               int32(FormV3),
				V4:               int32(FormV4),
				V5:               int32(FormV5),
				FirstkillT:       int32(FormFirstKillT),
				FirstkillCt:      int32(FormFirstKillCT),
				FirstdeathT:      int32(FormFirstDeathT),
				FirstdeathCt:     int32(FormFirstDeathCT),
			},
		},
	}, false)
}

// MatchVetoUpdateHandler Handler for /api/v1/match/{matchID}/vetoUpdate API. // TODO
func MatchVetoUpdateHandler(c *gin.Context) {
	/*log.Println("MatchVetoUpdateHandler")
	vars := mux.Vars(r)
	matchid, err := strconv.Atoi(vars["matchID"])
	if err != nil {
		http.Error(w, "matchid should be int", http.StatusBadRequest)
		return
	}

	FormMap := c.PostForm("map")
	FormTeamString := c.PostForm("teamString")
	FormPickOrVeto := c.PostForm("pick_or_veto")

	log.Printf("matchid : %d\n", matchid)
	log.Printf("FormMap : %s\n", FormMap)
	log.Printf("FormTeamString : %s\n", FormTeamString)
	log.Printf("FormPickOrVeto : %s\n", FormPickOrVeto)

	m := &db.MatchData{}
	db.SQLAccess.Gorm.First(&m,matchid)
	err = MatchAPICheck(m, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var TeamName string
	if FormTeamString == "team1" {
		team := &db.TeamData{}
		team.ID = m.Team1ID
		db.SQLAccess.Gorm.First(&team)
		TeamName = team.Name
	} else if FormTeamString == "team2" {
		team := &db.TeamData{}
		team.ID = m.Team2ID
		db.SQLAccess.Gorm.First(&team)
		TeamName = team.Name
	} else {
		TeamName = "Decider"
	}
	// veto = Veto.create(matchid, teamName, request.values.get('map'), request.values.get('pick_or_veto'))
	// TODO : Add Veto struct in db/models.go
	// Register to DB
	log.Printf("Confirmed Map Veto For %s on map %s\n", TeamName, FormMap)*/
	c.AbortWithError(http.StatusNotImplemented, fmt.Errorf("Not implemented yet"))
}

// MatchDemoUploadHandler Handler for /api/v1/match/{matchID}/map/{mapNumber}/demo API. // TODO
func MatchDemoUploadHandler(c *gin.Context) {
	/*log.Println("MatchDemoUploadHandler")
	vars := mux.Vars(r)
	matchid, err := strconv.Atoi(vars["matchID"])
	if err != nil {
		http.Error(w, "matchid should be int", http.StatusBadRequest)
		return
	}
	mapNumber, err := strconv.Atoi(vars["mapNumber"])
	if err != nil {
		http.Error(w, "mapNumber should be int", http.StatusBadRequest)
		return
	}

	DemoFile := c.PostForm("demoFile")

	m := &db.MatchData{}
	db.SQLAccess.Gorm.First(&m,matchid)
	err = MatchAPICheck(m, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mapstats := &db.MapStatsData{}
	mapstatsRecord := db.SQLAccess.Gorm.Where("match_id = ? AND map_number = ?", matchid, mapNumber).First(&mapstats)
	if !mapstatsRecord.RecordNotFound() {
		mapstatsRecord.Update("demoFile", DemoFile)
		log.Println("Made it through the demo post.")
		return
	}
	http.Error(w, "Failed to find map stats object", http.StatusBadRequest)
	*/
	c.AbortWithError(http.StatusNotImplemented, fmt.Errorf("Not implemented yet"))
}
