package v1

import (
	"database/sql"
	"github.com/FlowingSPDG/get5-web-go/server/src/util"
	"math"
	"time"
)

// UserData Struct for "user" table.
type APIUserData struct {
	ID      int                  `gorm:"primary_key;column:id;AUTO_INCREMENT" json:"id"`
	SteamID string               `gorm:"column:steam_id;unique" json:"steam_id"`
	Name    string               `gorm:"column:name" json:"name"`
	Admin   bool                 `gorm:"column:admin" json:"admin"`
	Servers []*APIGameServerData `gorm:"ForeignKey:UserID" json:"servers"`
	Teams   []*APITeamData       `gorm:"ForeignKey:UserID" json:"teams"`
	Matches []*APIMatchData      `gorm:"ForeignKey:UserID" json:"matches"`
}

// TableName declairation for GORM
func (u *APIUserData) TableName() string {
	return "user"
}

// GameServerData Struct for game_server table.
type APIGameServerData struct {
	ID           int    `gorm:"primary_key;column:id;AUTO_INCREMENT;NOT NULL" json:"id"`
	UserID       int    `gorm:"column:user_id;DEFAULT NULL" json:"user_id"`
	InUse        bool   `gorm:"column:in_use;DEFAULT NULL" json:"in_use"`
	IPString     string `gorm:"column:ip_string;DEFAULT NULL" json:"ip_string"`
	Port         int    `gorm:"column:port;DEFAULT NULL" json:"port"`
	Display      string `gorm:"column:display_name" json:"display_name"`
	PublicServer bool   `gorm:"column:public_server;DEFAULT NULL" json:"public_server"`
}

// TableName declairation for GORM
func (u *APIGameServerData) TableName() string {
	return "game_server"
}

// TeamData Struct for team table.
type APITeamData struct {
	ID          int      `gorm:"primary_key;column:id" json:"id"`
	UserID      int      `gorm:"column:user_id" json:"user_id"`
	Name        string   `gorm:"column:name" json:"name"`
	Tag         string   `gorm:"column:tag" json:"tag"`
	Flag        string   `gorm:"column:flag" json:"flag"`
	Logo        string   `gorm:"column:logo" json:"logo"`
	AuthsPickle []byte   `gorm:"column:auths" json:"-"`
	SteamIDs    []string `gorm:"-" json:"steamids"`
	PublicTeam  bool     `gorm:"column:public_team" json:"public_team"`

	User APIUserData `gorm:"ASSOCIATION_FOREIGNKEY:user_id" json:"-"`
}

// TableName declairation for GORM
func (t *APITeamData) TableName() string {
	return "team"
}

// GetPlayers Gets registered player's steamid64.
func (t *APITeamData) GetPlayers() ([]string, error) {
	auths, err := util.PickleToSteamID64s(t.AuthsPickle)
	if err != nil {
		return auths, err
	}
	return auths, nil
}

// APIMatchData Struct for match table.
type APIMatchData struct {
	ID            int          `gorm:"primary_key;column:id" json:"id"`
	UserID        int          `gorm:"column:user_id" json:"user_id"`
	Team1         APITeamData  `json:"team1"`
	Team2         APITeamData  `json:"team2"`
	Winner        int64        `gorm:"column:winner" json:"winner"`
	Cancelled     bool         `gorm:"column:cancelled" json:"cancelled"`
	StartTime     sql.NullTime `gorm:"column:start_time" json:"-"`
	StartTimeJSON time.Time    `json:"start_time"`
	EndTime       sql.NullTime `gorm:"column:end_time" json:"-"`
	EndTimeJSON   time.Time    `json:"end_time"`
	MaxMaps       int          `gorm:"column:max_maps" json:"max_maps"`
	Title         string       `gorm:"column:title" json:"title"`
	SkipVeto      bool         `gorm:"column:skip_veto" json:"skip_veto"`
	VetoMapPool   []string     `gorm:"column:veto_mappool" json:"veto_mappool"`
	Team1Score    int          `gorm:"column:team1_score" json:"team1_score"`
	Team2Score    int          `gorm:"column:team2_score" json:"team2_score"`
	Team1String   string       `gorm:"column:team1_string" json:"team1_string"`
	Team2String   string       `gorm:"column:team2_string" json:"team2_string"`
	Forfeit       bool         `gorm:"column:forfeit" json:"forfeit"`

	MapStats []APIMapStatsData `json:"map_stats"`
	Server   APIGameServerData `json:"server"`
	User     APIUserData       `json:"user"`

	Pending bool   `json:"pending"`
	Live    bool   `json:"live"`
	Status  string `json:"status"`
}

// TableName declairation for GORM
func (u *APIMatchData) TableName() string {
	return "match"
}

// APIMapStatsData MapStatsData struct for map_stats table.
type APIMapStatsData struct {
	ID            int          `gorm:"primary_key" gorm:"column:id" json:"id"`
	MatchID       int          `gorm:"column:match_id" gorm:"ForeignKey:match_id" json:"match_id"`
	MapNumber     int          `gorm:"column:map_number" json:"map_number"`
	MapName       string       `gorm:"column:map_name" json:"map_name"`
	StartTime     sql.NullTime `gorm:"column:start_time" json:"-"`
	StartTimeJSON time.Time    `json:"start_time"`
	EndTime       sql.NullTime `gorm:"column:end_time" json:"-"`
	EndTimeJSON   time.Time    `json:"end_time"`
	Winner        int          `gorm:"column:winner" json:"winner"`
	Team1Score    int          `gorm:"column:team1_score" json:"team1_score"`
	Team2Score    int          `gorm:"column:team2_score" json:"team2_score"`
}

// TableName declairation for GORM
func (u *APIMapStatsData) TableName() string {
	return "map_stats"
}

// APIPlayerStatsData Player stats data struct for player_stats table.
type APIPlayerStatsData struct {
	ID               int    `gorm:"primary_key;column:id" json:"id"`
	MatchID          int    `gorm:"column:match_id" json:"match_id"`
	MapID            int    `gorm:"column:map_id" json:"map_id"`
	TeamID           int    `gorm:"column:team_id" json:"team_id"`
	SteamID          string `gorm:"column:steam_id;unique" json:"steam_id"`
	Name             string `gorm:"column:name" json:"name"`
	Kills            int    `gorm:"column:kills" json:"kills"`
	Deaths           int    `gorm:"column:deaths" json:"deaths"`
	Roundsplayed     int    `gorm:"column:roundsplayed" json:"roundsplayed"`
	Assists          int    `gorm:"column:assists" json:"assists"`
	FlashbangAssists int    `gorm:"column:flashbang_assists" json:"flashbang_assists"`
	Teamkills        int    `gorm:"column:teamkills" json:"teamkills"`
	Suicides         int    `gorm:"column:suicides" json:"suicides"`
	HeadshotKills    int    `gorm:"column:headshot_kills" json:"headshot_kills"`
	Damage           int64  `gorm:"column:damage" json:"damage"`
	BombPlants       int    `gorm:"column:bomb_plants" json:"bomb_plants"`
	BombDefuses      int    `gorm:"column:bomb_defuses" json:"bomb_defuses"`
	V1               int    `gorm:"column:v1" json:"v1"`
	V2               int    `gorm:"column:v2" json:"v2"`
	V3               int    `gorm:"column:v3" json:"v3"`
	V4               int    `gorm:"column:v4" json:"v4"`
	V5               int    `gorm:"column:v5" json:"v5"`
	K1               int    `gorm:"column:k1" json:"k1"`
	K2               int    `gorm:"column:k2" json:"k2"`
	K3               int    `gorm:"column:k3" json:"k3"`
	K4               int    `gorm:"column:k4" json:"k4"`
	K5               int    `gorm:"column:k5" json:"k5"`
	FirstdeathCT     int    `gorm:"column:firstdeath_Ct" json:"firstdeath_Ct"`
	FirstdeathT      int    `gorm:"column:firstdeath_t" json:"firstdeath_t"`
	FirstkillCT      int    `gorm:"column:firstkill_ct" json:"firstkill_ct"`
	FirstkillT       int    `gorm:"column:firstkill_t" json:"firstkill_t"`

	Rating float64 `json:"rating"`
	KDR    float64 `json:"kdr"`
	HSP    float64 `json:"hsp"`
	ADR    float64 `json:"adr"`
	FPR    float64 `json:"fpr"`
}

// TableName declairation for GORM
func (p *APIPlayerStatsData) TableName() string {
	return "player_stats"
}

// GetRating Get player's rating. Average datas are static tho.
func (p *APIPlayerStatsData) GetRating() float64 { // Rating value can be more accurate??
	var AverageKPR float64 = 0.679
	var AverageSPR float64 = 0.317
	var AverageRMK float64 = 1.277
	var KillRating float64 = float64(p.Kills) / float64(p.Roundsplayed) / AverageKPR
	var SurvivalRating float64 = float64(p.Roundsplayed-p.Deaths) / float64(p.Roundsplayed) / float64(AverageSPR)
	var killcount float64 = float64(p.K1 + 4*p.K2 + 9*p.K3 + 16*p.K4 + 25*p.K5)
	var RoundsWithMultipleKillsRating float64 = killcount / float64(p.Roundsplayed) / float64(AverageRMK)
	var rating float64 = (KillRating + 0.7*SurvivalRating + RoundsWithMultipleKillsRating) / 2.7
	if math.IsNaN(rating) {
		return 0.0
	}
	return util.Round(rating, 2)
}

// GetKDR Returns player's KDR(Kill/Deaths Ratio).
func (p *APIPlayerStatsData) GetKDR() float64 {
	if p.Deaths == 0 {
		return float64(p.Kills)
	}
	return util.Round(float64(p.Kills)/float64(p.Deaths), 2)
}

// GetHSP Returns player's HSP(HeadShot Percentage).
func (p *APIPlayerStatsData) GetHSP() float64 {
	if p.Kills == 0 {
		return 0
	}
	return util.Round(float64(p.HeadshotKills)/float64(p.Kills)*100, 2)
}

// GetADR Returns player's ADR(Average Damage per Round).
func (p *APIPlayerStatsData) GetADR() float64 {
	if p.Roundsplayed == 0 {
		return 0.0
	}
	return util.Round(float64(p.Damage)/float64(p.Roundsplayed), 2)
}

// GetFPR Returns player's FPR(Frags Per Round).
func (p *APIPlayerStatsData) GetFPR() float64 {
	if p.Roundsplayed == 0 {
		return 0.0
	}
	return util.Round(float64(p.Kills)/float64(p.Roundsplayed), 2)
}
