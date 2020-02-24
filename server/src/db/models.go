package db

import (
	"database/sql"
	"fmt"
	"github.com/FlowingSPDG/get5-web-go/server/src/cfg"
	"github.com/FlowingSPDG/get5-web-go/server/src/util"
	"log"
	"net"
	"time"

	"strings"
	//"github.com/Philipp15b/go-steam"
	//"github.com/FlowingSPDG/go-steamapi"
	"github.com/Acidic9/steam"
	"github.com/jinzhu/gorm"
	gosteam "github.com/kidoman/go-steam"

	//"github.com/solovev/steam_go"
	//_ "html/template"
	"math"
	"strconv"
)

// UserData Struct for "user" table.
type UserData struct {
	ID      int    `gorm:"primary_key;column:id;AUTO_INCREMENT"`
	SteamID string `gorm:"column:steam_id;unique"`
	Name    string `gorm:"column:name"`
	Admin   bool   `gorm:"column:admin"`

	Servers []*GameServerData `gorm:"ForeignKey:UserID"`
	Teams   []*TeamData       `gorm:"ForeignKey:UserID"`
	Matches []*MatchData      `gorm:"ForeignKey:UserID"`
}

// TableName declairation for GORM
func (u *UserData) TableName() string {
	return "user"
}

// GetOrCreateBySteamID Get or Register Userdata by their steamid into DB.
func (u *UserData) GetOrCreateBySteamID() (*UserData, bool, error) {
	return u.GetOrCreate()
}

// GetOrCreate Get or Register Userdata by their steamid into DB.
func (u *UserData) GetOrCreateByID(userid int) (*UserData, bool, error) { // userdata, exist,error
	SQLUserData := UserData{}
	exist := false

	record := SQLAccess.Gorm.FirstOrCreate(userid, &SQLUserData)
	if record.RecordNotFound() {
		exist = false
		log.Println("USER NOT EXIST!")
		log.Println("CREATING USER")
		steamid64, err := strconv.Atoi(u.SteamID)
		if err != nil {
			return u, exist, err
		}
		SQLUserData.Name, err = GetSteamName(uint64(steamid64))
		if err != nil {
			return u, exist, err
		}
		SQLAccess.Gorm.Create(&SQLUserData)
		SQLAccess.Gorm.Where("steam_id = ?", u.SteamID).First(u)
	} else {
		log.Println("USER EXIST")
		exist = true
		// log.Println(SQLUserData)
		u.Name = SQLUserData.Name
		u.ID = SQLUserData.ID
		u.Admin = SQLUserData.Admin
		u.SteamID = SQLUserData.SteamID
	}
	return u, exist, nil
}

// GetOrCreate Get or Register Userdata by their steamid into DB.
func (u *UserData) GetOrCreate() (*UserData, bool, error) { // userdata, exist,error
	SQLUserData := UserData{}
	exist := false

	record := SQLAccess.Gorm.Where("steam_id = ?", u.SteamID).First(&SQLUserData)
	if record.RecordNotFound() {
		exist = false
		log.Println("USER NOT EXIST!")
		log.Println("CREATING USER")
		steamid64, err := strconv.Atoi(u.SteamID)
		if err != nil {
			return u, exist, err
		}
		SQLUserData.Name, err = GetSteamName(uint64(steamid64))
		if err != nil {
			return u, exist, err
		}
		SQLAccess.Gorm.Create(&SQLUserData)
		SQLAccess.Gorm.Where("steam_id = ?", u.SteamID).First(u)
	} else {
		log.Println("USER EXIST")
		exist = true
		// log.Println(SQLUserData)
		u.Name = SQLUserData.Name
		u.ID = SQLUserData.ID
		u.Admin = SQLUserData.Admin
		u.SteamID = SQLUserData.SteamID
	}
	return u, exist, nil
}

// GetURL Get user page URL
func (u *UserData) GetURL() string {
	return fmt.Sprintf("http://%s/user/%d", config.Cnf.HOST, u.ID)
}

// GetSteamURL Get user's steam page URL by their steamid64
func (u *UserData) GetSteamURL() string {
	return "http://steamcommunity.com/profiles/" + u.SteamID
}

// GetRecentMatches Gets match history
func (u *UserData) GetRecentMatches(limit int) []*MatchData {
	SQLAccess.Gorm.Model(u).Related(&u.Matches, "Matches")
	return u.Matches
}

// GetTeams Get teams which is owened by user
func (u *UserData) GetTeams(limit int) []*TeamData {
	SQLAccess.Gorm.Model(u).Related(&u.Teams, "Teams")
	return u.Teams
}

// GameServerData Struct for game_server table.
type GameServerData struct {
	ID           int    `gorm:"primary_key;column:id;AUTO_INCREMENT;NOT NULL" json:"id"`
	UserID       int    `gorm:"column:user_id;DEFAULT NULL" json:"user_id"`
	InUse        bool   `gorm:"column:in_use;DEFAULT NULL" json:"in_use"`
	IPString     string `gorm:"column:ip_string;DEFAULT NULL" json:"ip_string"`
	Port         int    `gorm:"column:port;DEFAULT NULL" json:"port"`
	RconPassword string `gorm:"column:rcon_password;DEFAULT NULL" json:"rcon_password"`
	DisplayName  string `gorm:"column:display_name;DEFAULT NULL" json:"display_name"`
	PublicServer bool   `gorm:"column:public_server;DEFAULT NULL" json:"public_server"`
}

// TableName declairation for GORM
func (g *GameServerData) TableName() string {
	return "game_server"
}

// Create Register GameServer into DB.
func (g *GameServerData) Create(userid int, displayname string, ipstring string, port int, rconpassword string, publicserver bool) (*GameServerData, error) {
	if ipstring == "" || rconpassword == "" {
		return nil, fmt.Errorf("IPaddress or RCON password is empty")
	}
	var servernum int
	user := UserData{}
	// returns error if user is not valid.
	rec := SQLAccess.Gorm.First(&user, userid)
	/*if rec.Error != nil {
		return nil, rec.Error
	}*/
	if !rec.RecordNotFound() {
		SQLAccess.Gorm.Model(&user).Related(&user.Servers, "Servers").Count(&servernum)
		if uint16(servernum) > config.Cnf.UserMaxResources.Servers {
			return nil, fmt.Errorf("Max server limit exceeded! OWNED[%d] / LIMIT:[%d]", servernum, config.Cnf.UserMaxResources.Servers)
		}
	}

	g.UserID = userid
	g.DisplayName = displayname
	g.IPString = ipstring
	g.Port = port
	g.RconPassword = rconpassword
	g.PublicServer = publicserver

	_, err := util.CheckServerAvailability(g.IPString, g.Port, g.RconPassword)
	if err != nil {
		return nil, err
	}

	result := SQLAccess.Gorm.Create(&g)
	errors := result.GetErrors()
	log.Printf("Errors len : %d\n", len(errors))
	if len(errors) >= 1 {
		return nil, errors[0]
	}
	return g, result.Error
}

// CanEdit Check if server is editable for user or not.
func (g *GameServerData) CanEdit(userid int) bool {
	if g.UserID == 0 {
		SQLAccess.Gorm.First(&g, g.ID)
	}
	if userid == 0 {
		return false
	}
	return g.UserID == userid
}

// CanDelete Check if server is deletable for user or not.
func (g *GameServerData) CanDelete(userid int) bool {
	return g.CanEdit(userid)
}

// Edit Edits GameServer information.
func (g *GameServerData) Edit() (*GameServerData, error) {
	if g.ID == 0 {
		return nil, fmt.Errorf("ID not valid")
	}
	Server := GameServerData{}
	rec := SQLAccess.Gorm.First(&Server, g.ID)
	if rec.RecordNotFound() {
		return nil, fmt.Errorf("Server not found")
	}
	SQLAccess.Gorm.Model(&Server).Update(&g)
	SQLAccess.Gorm.Save(&g)
	return g, nil
}

// Delete Deletes GameServer information.
func (g *GameServerData) Delete() error {
	if g.ID == 0 {
		return fmt.Errorf("ID not valid")
	}
	rec := SQLAccess.Gorm.First(&g, g.ID)
	if rec.RecordNotFound() {
		return fmt.Errorf("Server not found")
	}
	SQLAccess.Gorm.Delete(&g)
	return nil
}

// SendRcon Sends Remote-Commands to specific IP SRCDS.
func (g *GameServerData) SendRcon(cmd string) (string, error) {
	if !checkIP(g.IPString) {
		return "", fmt.Errorf("Specified IP is not valid")
	}
	o := &gosteam.ConnectOptions{RCONPassword: g.RconPassword}
	rcon, err := gosteam.Connect(g.GetHostPort(), o)
	if err != nil {
		return "", err
	}
	defer rcon.Close()

	resp, err := rcon.Send(cmd)
	if err != nil {
		return "", err
	}
	return resp, nil
}

// GetHostPort gets gameserver addr and port. e.g. localhost:27015
func (g *GameServerData) GetHostPort() string {
	return fmt.Sprintf("%s:%d", g.IPString, g.Port)
}

// GetDisplay Returns "DisplayName" if its not empty. otherwise it returns address and port.
func (g *GameServerData) GetDisplay() string {
	if g.DisplayName != "" {
		return fmt.Sprintf("%s (%s)", g.DisplayName, g.GetHostPort())
	}
	return g.GetHostPort()
}

// TeamData Struct for team table.
type TeamData struct {
	ID           int               `gorm:"primary_key;column:id;AUTO_INCREMENT" json:"id"`
	UserID       int               `gorm:"column:user_id" json:"user_id"`
	Name         string            `gorm:"column:name" json:"name"`
	Tag          string            `gorm:"column:tag" json:"tag"`
	Flag         string            `gorm:"column:flag" json:"flag"`
	Logo         string            `gorm:"column:logo" json:"logo"`
	AuthsPickle  []byte            `gorm:"column:auths" json:"-"`
	SteamIDs     string            `gorm:"column:steamids" json:"-"`
	SteamIDsJSON []string          `gorm:"-" json:"auths"`
	Players      []PlayerStatsData `gorm:"-" json:"-"`
	PublicTeam   bool              `gorm:"column:public_team" json:"public_team"`
	MixTeam      bool              `gorm:"column:mix_team" json:"mix_team"`
}

// TableName declairation for GORM
func (t *TeamData) TableName() string {
	return "team"
}

// Create Register Team information into DB.
func (t *TeamData) Create(userid int, name string, tag string, flag string, logo string, auths []string, publicteam bool, mixteam bool) (*TeamData, error) {
	if name == "" {
		return nil, fmt.Errorf("Team name cannot be empty")
	}
	var teamnum int
	user := UserData{}
	rec := SQLAccess.Gorm.First(&user, userid)
	/*if rec.Error != nil {
		return nil, rec.Error
	}*/
	if !rec.RecordNotFound() {
		SQLAccess.Gorm.Model(&user).Related(&user.Teams, "Teams").Count(&teamnum)
		if uint16(teamnum) > config.Cnf.UserMaxResources.Teams {
			return nil, fmt.Errorf("Max teams limit exceeded! OWNED[%d] / LIMIT:[%d]", teamnum, config.Cnf.UserMaxResources.Teams)
		}
	}
	t.ID = 0
	t.UserID = userid
	t.Name = name
	t.Tag = tag
	t.Flag = flag
	t.Logo = logo
	t.PublicTeam = publicteam
	t.MixTeam = mixteam
	var err error
	for i := 0; i < len(auths); i++ {
		auths[i], err = util.AuthToSteamID64(auths[i])
		if err != nil {
			return nil, err
		}
	}
	t.SteamIDs = strings.Join(auths, ",")
	t.SteamIDsJSON = auths
	t.AuthsPickle, err = util.SteamID64sToPickle(auths)
	if err != nil {
		return nil, err
	}
	rec = SQLAccess.Gorm.Create(&t)
	if rec.Error != nil {
		return nil, rec.Error
	}
	return t, nil
}

// Edit Edits TeamData information.
func (t *TeamData) Edit() (*TeamData, error) {
	if t.ID == 0 {
		return nil, fmt.Errorf("ID not valid")
	}
	Team := TeamData{}
	rec := SQLAccess.Gorm.First(&Team, t.ID)
	if rec.RecordNotFound() {
		return nil, fmt.Errorf("Team not found")
	}

	var err error
	for i := 0; i < len(t.SteamIDsJSON); i++ {
		if t.SteamIDsJSON[i] != "" {
			t.SteamIDsJSON[i], err = util.AuthToSteamID64(t.SteamIDsJSON[i])
			if err != nil {
				return nil, err
			}
		}
	}
	t.AuthsPickle, err = util.SteamID64sToPickle(t.SteamIDsJSON)
	t.SteamIDs = strings.Join(t.SteamIDsJSON, ",")
	SQLAccess.Gorm.Model(&Team).Update(&t)
	SQLAccess.Gorm.Save(&t)
	return t, nil
}

// Delete Deletes TeamData information.
func (t *TeamData) Delete() error {
	if t.ID == 0 {
		return fmt.Errorf("ID not valid")
	}
	rec := SQLAccess.Gorm.Delete(&t, t.ID)
	if rec.RecordNotFound() {
		return fmt.Errorf("Team not found")
	}
	if rec.Error != nil {
		return rec.Error
	}
	return nil
}

// CanEdit Check if server is editable for user or not.
func (t *TeamData) CanEdit(userid int) bool {
	if t.UserID == 0 {
		SQLAccess.Gorm.First(&t, t.ID)
	}
	if userid == 0 {
		return false
	} else if t.UserID == userid {
		return true
	}
	return false
}

// CanDelete Check if server is deletable for user or not.
func (t *TeamData) CanDelete(userid int) bool {
	if t.CanEdit(userid) == false {
		return false
	}
	// return len(t.GetRecentMatches(10)) == 0 // ?
	return true
}

// GetPlayers Gets registered player's steamid64.
func (t *TeamData) GetPlayers() (*[]string, error) {
	auths, err := util.PickleToSteamID64s(t.AuthsPickle)
	var results = make([]string, 0, len(auths))
	if err != nil {
		return &results, err
	}
	for i := 0; i < len(auths); i++ {
		results = append(results, auths[i])
	}
	/*for i := 0; i < len(auths); i++ {
		steamid64, err := strconv.Atoi(auths[i])
		if err != nil {
			return results, err
		}
		playername, err := GetSteamName(uint64(steamid64))
		if err != nil {
			return results, err
		}
		//results = append(results, GetPlayerData{Auth: auths[i], Name: playername})
		results = append(results, GetPlayerData{Auth: auths[i]})
	}*/
	t.SteamIDsJSON = results
	t.SteamIDs = strings.Join(t.SteamIDsJSON, ",")

	return &results, nil
}

/*
func (t *TeamData) GetPlayers() ([]PlayerStatsData, error) {
	reader := bytes.NewReader(t.AuthsPickle)
	var results []string
	err := stalecucumber.UnpackInto(&results).From(stalecucumber.Unpickle(reader))
	if err != nil {
		return t.Players, err
	}
	// it won't return datas in case team never played game...
		t.Players = []PlayerStatsData{}
		for i := 0; i < len(results); i++ {
			if results[i] != "" {
				p := PlayerStatsData{}
				SQLAccess.Gorm.Where("steam_id = ?", results[i]).First(&p)
				log.Println(p)
				if err != nil {
					return t.Players, err
				}
				if p.Steam_id != "" {
					t.Players = append(t.Players, p)
				}
			}
		}
	//SQLAccess.Gorm.Where("team_id = ?", t.ID).Find(&t.Players) // N+1 issue
	return t.Players, nil
}
*/

// GetRecentMatches Gets team match history.
func (t *TeamData) GetRecentMatches(limit int) []*MatchData {
	var matches []*MatchData
	if t.PublicTeam == true {
		SQLAccess.Gorm.Where("team1_id = ?", t.ID).Or("team2_id = ?", t.ID).Not("start_time = null AND cancelled = true").Limit(limit).Find(&matches)
	} else {
		var owner UserData
		SQLAccess.Gorm.First(&owner, t.UserID)
		SQLAccess.Gorm.Model(&owner).Related(&owner.Matches, "Matches").Limit(limit)
		matches = owner.Matches
	}
	return matches
}

// GetVSMatchResult Returns Match result as string for gorazor template.
func (t *TeamData) GetVSMatchResult(matchid int) (string, error) {
	var otherteam TeamData
	myscore := 0
	otherteamscore := 0
	var match MatchData
	SQLAccess.Gorm.First(&match, matchid)
	if int(match.Team1ID) == t.ID {
		myscore = match.Team1Score
		otherteamscore = match.Team2Score
		SQLAccess.Gorm.First(&otherteam, match.Team2ID)
	} else {
		myscore = match.Team2Score
		otherteamscore = match.Team1Score
		SQLAccess.Gorm.First(&otherteam, match.Team2ID)
	}

	// for a bo1 replace series score with the map score...
	if match.MaxMaps == 1 {
		_, err := match.GetMapStat()
		if err != nil {
			return "", err
		}
		if len(match.MapStats) <= 0 {
			return fmt.Sprintf("Pending, vs %s", otherteam.Name), nil // maybe add <a> tag for otherteam.Name ?
		}
		mapstat := match.MapStats[0]
		if int(match.Team1ID) == t.ID {
			myscore = mapstat.Team1Score
			otherteamscore = mapstat.Team2Score
		} else {
			myscore = mapstat.Team2Score
			otherteamscore = mapstat.Team1Score
		}
	}
	if match.Live() == true {
		return fmt.Sprintf("Live, %d:%d vs %s", myscore, otherteamscore, otherteam.Name), nil // maybe add <a> tag for otherteam.Name ?
	}
	if myscore < otherteamscore {
		return fmt.Sprintf("Lost %d:%d vs %s", myscore, otherteamscore, otherteam.Name), nil
	} else if myscore > otherteamscore {
		return fmt.Sprintf("Won %d:%d vs %s", otherteamscore, myscore, otherteam.Name), nil
	} else {
		return fmt.Sprintf("Tied %d:%d vs %s", otherteamscore, myscore, otherteam.Name), nil
	}

}

// GetURL Get URL of team page.
func (t *TeamData) GetURL() string {
	return fmt.Sprintf("http://%s/team/%d", config.Cnf.HOST, t.ID)
}

// GetNameURLHtml Get team page and name as a-tag. for gorazor template
func (t *TeamData) GetNameURLHtml() string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, t.GetURL(), t.Name)
}

// MatchConfig Match configration for get5 api, based on https://github.com/splewis/get5/blob/master/configs/get5/example_match.json and https://github.com/splewis/get5/blob/3f793ceb3736d78ba6a92f42631d91cb52f0beb4/scripting/get5/matchconfig.sp#L435
type MatchConfig struct {
	MatchID              string `json:"matchid"` // NOT int
	Scrim                bool   `json:"scrim,omitempty"`
	MatchTitle           string `json:"match_title,omitempty"`
	PlayersPerTeam       int    `json:"players_per_team,omitempty"`
	MinPlayersToReady    int    `json:"min_players_to_ready,omitempty"`
	MinSPectatorsToReady int    `json:"min_spectators_to_ready,omitempty"`
	SkipVeto             bool   `json:"skip_veto,omitempty"`
	// bo2_series and maps_to_win are deprecated on get5 plugin side. Use num_maps insted
	// Bo2Series 			   bool   `json:"bo2_series"`
	// MapsToWin 			   int    `json:"maps_to_win"`
	NumMaps                    int    `json:"num_maps,omitempty"`
	VetoFirst                  string `json:"veto_first,omitempty"` // team1 || team2
	SideType                   string `json:"side_type,omitempty"`  // standard and ?
	FavoredTeamPercentageText  string `json:"favored_percentage_text,omitempty"`
	FavoredTeamPercentageTeam1 int    `json:"favored_percentage_team1,omitempty"`
	Spectators                 struct {
		Name    string   `json:"name"`
		Players []string `json:"players"`
	} `json:"spectators,omitempty"`
	Team1 struct {
		Flag    string   `json:"flag"`
		Name    string   `json:"name"`
		Tag     string   `json:"tag"`
		Players []string `json:"players"`
	} `json:"team1"`
	Team2 struct {
		Flag    string   `json:"flag"`
		Name    string   `json:"name"`
		Tag     string   `json:"tag"`
		Players []string `json:"players"`
	} `json:"team2"`

	Maplist  []string `json:"maplist"`
	MapSides []string `json:"map_sides,omitempty"`

	Cvars map[string]string `json:"cvars,omitempty"`
}

// MatchData Struct for match table.
type MatchData struct {
	// Original columns...
	ID              int           `gorm:"primary_key;column:id;AUTO_INCREMENT" json:"id"`
	UserID          int           `gorm:"column:user_id" json:"user_id"`
	ServerID        int           `gorm:"column:server_id" json:"server_id"`
	Team1ID         int           `gorm:"column:team1_id" json:"team1_id"`
	Team2ID         int           `gorm:"column:team2_id" json:"team2_id"`
	Winner          sql.NullInt32 `gorm:"column:winner" json:"winner"`
	Cancelled       bool          `gorm:"column:cancelled" json:"cancelled"`
	StartTime       sql.NullTime  `gorm:"column:start_time" json:"start_time"`
	EndTime         sql.NullTime  `gorm:"column:end_time" json:"end_time"`
	MaxMaps         int           `gorm:"column:max_maps" json:"max_maps"`
	Title           string        `gorm:"column:title" json:"title"`
	SkipVeto        bool          `gorm:"column:skip_veto" json:"skip_veto"`
	APIKey          string        `gorm:"column:api_key" json:"-"`
	VetoMapPool     string        `gorm:"column:veto_mappool" json:"-"`
	VetoMapPoolJSON []string      `gorm:"-" json:"veto_mappool"`
	Team1Score      int           `gorm:"column:team1_score" json:"team1_score"`
	Team2Score      int           `gorm:"column:team2_score" json:"team2_score"`
	Team1String     string        `gorm:"column:team1_string" json:"team1_string"`
	Team2String     string        `gorm:"column:team2_string" json:"team2_string"`
	Forfeit         bool          `gorm:"column:forfeit" json:"forfeit"`
	PluginVersion   string        `gorm:"column:plugin_version" json:"-"`

	// get5-web-go columns...
	CvarsJSON map[string]string `gorm:"-" json:"cvars"`
	Cvars     string            `gorm:"column:cvars" json:"-"`
	SideType  string            `gorm:"column:side_type" json:"side_type"`
	IsPug     bool              `gorm:"column:is_pug" json:"-"`

	// If its pug
	SteamIDs []string `gorm:"-" json:"steamids"`

	MapStats []MapStatsData `gorm:"ForeignKey:MatchID" json:"-"`
	Server   GameServerData `json:"-"`
}

// TableName declairation for GORM
func (m *MatchData) TableName() string {
	return "match"
}

// Get Get myself
func (m *MatchData) Get(id int) (*MatchData, error) {
	rec := SQLAccess.Gorm.First(&m, id)
	if rec.RecordNotFound() {
		return nil, fmt.Errorf("Match not found")
	}
	return m, nil
}

// Create Register Match information into DB.
func (m *MatchData) Create(userid int, team1id int, team2id int, team1string string, team2string string, maxmaps int, skipveto bool, title string, vetomappool []string, serverid int, cvars map[string]string, sidetype string, pug bool) (*MatchData, error) {
	user := UserData{}

	var matchnum int
	rec := SQLAccess.Gorm.Model(&MatchData{}).Where("user_id = ?", userid).Count(&matchnum)
	/*if rec.Error != nil {
		return nil, rec.Error
	}*/
	if !rec.RecordNotFound() {
		SQLAccess.Gorm.Model(&user).Related(&user.Matches, "Matches").Count(&matchnum)
		if uint16(matchnum) > config.Cnf.UserMaxResources.Matches {
			return nil, fmt.Errorf("Max match limit exceeded! OWNED[%d] / LIMIT:[%d]", matchnum, config.Cnf.UserMaxResources.Matches)
		}
	}

	rec = SQLAccess.Gorm.First(&user, userid)
	if team1id == 0 || team2id == 0 || serverid == 0 {
		return nil, fmt.Errorf("TeamID or ServerID is empty")
	}
	server := GameServerData{}
	rec = SQLAccess.Gorm.First(&server, serverid)
	if rec.Error != nil {
		return nil, rec.Error
	}
	// returns error if user wasnt owned server,or not an admin.
	if userid != server.UserID && !user.Admin && !server.PublicServer {
		return nil, fmt.Errorf("This is not your server")
	}

	get5res, err := util.CheckServerAvailability(server.IPString, server.Port, server.RconPassword) // Returns error if SRCDS is not available
	if err != nil {
		return nil, err
	}

	MatchOnServer := MatchData{}
	rec = SQLAccess.Gorm.Where("end_time = NULL AND server_id = ? AND cancelled = FALSE", serverid).First(&MatchOnServer)
	if !rec.RecordNotFound() {
		return nil, rec.Error
	}

	m.UserID = userid
	m.ServerID = serverid
	m.GetServer()
	m.Team1ID = team1id
	m.Team2ID = team2id
	m.MaxMaps = maxmaps
	m.Title = title
	m.SkipVeto = skipveto
	m.VetoMapPool = strings.Join(vetomappool, " ")
	m.Team1String = team1string
	m.Team2String = team2string
	if get5res.PluginVersion == "" {
		get5res.PluginVersion = "unknown"
	}

	for k, v := range cvars {
		m.Cvars += fmt.Sprintf("%s %s", k, v) // [ "hostname" = "WASD","tv_enable" = "1" ]
		m.Cvars += "\n"
	}
	m.Cvars = strings.TrimRight(m.Cvars, "\n") // trim last "\n".
	log.Printf("m.Cvars : %v\n", m.Cvars)

	m.SideType = sidetype
	m.IsPug = pug

	MatchServerUpdate := m.Server
	rec = SQLAccess.Gorm.First(&MatchServerUpdate)
	if rec.Error != nil {
		return nil, rec.Error
	}
	MatchServerUpdate.InUse = true
	rec = SQLAccess.Gorm.Model(&m.Server).Update(&MatchServerUpdate)
	if rec.Error != nil {
		return nil, rec.Error
	}
	rec = SQLAccess.Gorm.Save(&MatchServerUpdate)
	if rec.Error != nil {
		return nil, rec.Error
	}

	m.PluginVersion = get5res.PluginVersion
	m.APIKey = util.RandString(24)
	rec = SQLAccess.Gorm.Create(&m)
	log.Printf("CREATE RESULT : %v\n", rec)
	errors := rec.GetErrors()
	if len(errors) != 0 {
		return nil, errors[0]
	}
	if rec.Error != nil {
		return nil, err
	}
	err = m.SendToServer()
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetStatusString Get match status as string.
func (m *MatchData) GetStatusString(ShowWinner bool) (string, error) {
	if m.Pending() {
		return "Pending", nil
	} else if m.Live() {
		teams1core, team2score, err := m.GetCurrentScore(SQLAccess.Gorm)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Live, %d:%d", teams1core, team2score), nil
	} else if m.Finished() {
		teams1core, team2score, err := m.GetCurrentScore(SQLAccess.Gorm)
		if err != nil {
			return "", err
		}
		minscore := math.Min(float64(teams1core), float64(team2score))
		maxscore := math.Max(float64(teams1core), float64(team2score))
		ScoreString := fmt.Sprintf("%d:%d", int(maxscore), int(minscore))
		winner := int(m.Winner.Int32)
		if !ShowWinner {
			return "Finished", nil
		} else if winner == m.Team1ID {
			team1, err := m.GetTeam1()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Won %s by %s", ScoreString, team1.Name), nil
		} else if winner == m.Team2ID {
			team2, err := m.GetTeam2()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Won %s by %s", ScoreString, team2.Name), nil
		} else {
			return fmt.Sprintf("Tied %s", ScoreString), nil
		}
	} else {
		return "Cancelled", nil
	}
}

// GetVSString Get Match VS information as string. for gorazor template
func (m *MatchData) GetVSString() (string, error) {
	team1, err := m.GetTeam1()
	team2, err := m.GetTeam2()
	if err != nil {
		return "", err
	}
	team1score, team2score, err := m.GetCurrentScore(SQLAccess.Gorm)
	if err != nil {
		return "", err
	}
	str := fmt.Sprintf("%s VS %s (%d:%d)", team1.GetNameURLHtml(), team2.GetNameURLHtml(), team1score, team2score)
	return str, nil
}

// Finalized Returns true if match is finished or cancelled
func (m *MatchData) Finalized() bool {
	return m.Cancelled || m.Finished()
}

// Pending Returns true if match is not started and not cancelled
func (m *MatchData) Pending() bool {
	return !m.StartTime.Valid && !m.Cancelled
}

// Finished Returns true if match is ended and not cancelled
func (m *MatchData) Finished() bool {
	return m.EndTime.Valid && !m.Cancelled
}

// Live Retursn true if match is in-progress
func (m *MatchData) Live() bool {
	return m.StartTime.Valid && !m.EndTime.Valid && !m.Cancelled
}

// GetServer Get match server ID as GameServerData
func (m *MatchData) GetServer() (*GameServerData, error) {
	rec := SQLAccess.Gorm.Model(m).Related(&m.Server, "Servers")
	if rec.RecordNotFound() {
		return nil, fmt.Errorf("Server not found")
	}
	return &m.Server, nil
}

// GetCurrentScore Returns current match score. returns map-score if match is BO1.
func (m *MatchData) GetCurrentScore(g *gorm.DB) (int, int, error) {
	if m.ID == 0 {
		return 0, 0, fmt.Errorf("Match ID invalid")
	}
	m.GetMapStat()
	if m.MaxMaps == 1 {
		if len(m.MapStats) == 0 { // check ok?
			return 0, 0, nil
		}
		return m.MapStats[0].Team1Score, m.MapStats[0].Team2Score, nil
	}
	return m.Team1Score, m.Team2Score, nil
}

// GetTeam1 Get Team1 as "TeamData" struct.
func (m *MatchData) GetTeam1() (TeamData, error) {
	var Team = TeamData{}
	var STeam TeamData
	SQLAccess.Gorm.First(&STeam, m.Team1ID)
	Team.ID = STeam.ID
	Team.Name = STeam.Name
	Team.Tag = STeam.Tag
	Team.Flag = STeam.Flag
	Team.Logo = STeam.Logo
	Team.AuthsPickle = STeam.AuthsPickle
	Team.PublicTeam = STeam.PublicTeam
	var err error
	Team.SteamIDsJSON, err = util.PickleToSteamID64s(STeam.AuthsPickle)
	if err != nil {
		return Team, err
	}
	return Team, nil
}

// GetTeam2 Get Team2 as "TeamData" struct.
func (m *MatchData) GetTeam2() (TeamData, error) {
	var Team = TeamData{}
	var STeam TeamData
	SQLAccess.Gorm.First(&STeam, m.Team2ID)
	Team.ID = STeam.ID
	Team.Name = STeam.Name
	Team.Tag = STeam.Tag
	Team.Flag = STeam.Flag
	Team.Logo = STeam.Logo
	Team.AuthsPickle = STeam.AuthsPickle
	Team.PublicTeam = STeam.PublicTeam
	var err error
	Team.SteamIDsJSON, err = util.PickleToSteamID64s(STeam.AuthsPickle)
	if err != nil {
		return Team, err
	}
	return Team, nil
}

/*
// GetUser Get Match owner as "UserData" struct.
func (m *MatchData) GetUser() UserData {
	SQLAccess.Gorm.First(&m.User, m.UserID)
	return m.User
}
*/

// GetWinner Get Winner team as "TeamData" struct.
func (m *MatchData) GetWinner() (TeamData, error) {
	var Empty TeamData
	if m.Team1Score > m.Team2Score {
		winner, err := m.GetTeam1()
		if err != nil {
			return Empty, err
		}
		return winner, nil
	} else if m.Team2Score > m.Team1Score {
		winner, err := m.GetTeam2()
		if err != nil {
			return Empty, err
		}
		return winner, nil
	}
	return Empty, nil
}

// GetLoser Get Loser team as "TeamData" struct.
func (m *MatchData) GetLoser() (TeamData, error) {
	var Empty TeamData
	if m.Team1Score > m.Team2Score {
		loser, err := m.GetTeam2()
		if err != nil {
			return Empty, err
		}
		return loser, nil
	} else if m.Team2Score > m.Team1Score {
		loser, err := m.GetTeam1()
		if err != nil {
			return Empty, err
		}
		return loser, nil
	}
	return Empty, nil
}

// SendToServer Send match config to server
func (m *MatchData) SendToServer() error {
	if m.ServerID == 0 || m.Server.ID == 0 {
		return fmt.Errorf("Server not found")
	}

	m.Server.SendRcon("log on")
	res, err := m.Server.SendRcon(fmt.Sprintf("get5_loadmatch_url %s/api/v1/match/%d/config", config.Cnf.HOST, m.ID))
	res, err = m.Server.SendRcon(fmt.Sprintf("get5_web_api_key %s", m.APIKey))
	res, err = m.Server.SendRcon(fmt.Sprintf("logaddress_add_http \"http://%s/api/v1/match/%d/csgolog/%s\"", config.Cnf.HOST, m.ID, m.APIKey))
	if err != nil || res != "" {
		return err
	}
	SQLAccess.Gorm.First(&m.Server)
	SQLAccess.Gorm.Model(&m.Server).Update("in_use", true)
	return nil
}

// BuildMatchDict Builds match JSON data.
func (m *MatchData) BuildMatchDict() (*MatchConfig, error) {
	SQLAccess.Gorm.First(&m, m.ID)
	m.VetoMapPoolJSON = strings.Split(m.VetoMapPool, " ")
	team1, err := m.GetTeam1()
	team2, err := m.GetTeam2()
	if err != nil {
		return &MatchConfig{}, err
	}
	log.Printf("team1 : %v\n", team1)
	log.Printf("team2 : %v\n", team2)
	var cfg = MatchConfig{
		MatchID: strconv.Itoa(m.ID),
		//Scrim:false,
		MatchTitle:        m.Title,
		PlayersPerTeam:    5, // 0 broke Veto commencing section // not 1 // original get5-web did not specify this value...
		MinPlayersToReady: 1, // Minimum # of players a team must have to ready // original get5-web did not specify this value...
		// MinSpectatorsToReady: // How many spectators must be ready to begin.
		SkipVeto: m.SkipVeto, // If set to 1, the maps will be preset using the first maps in the maplist below.
		NumMaps:  m.MaxMaps,  // Must be an odd number or 2. 1->Bo1, 2->Bo2, 3->Bo3, etc.
		// VetoFirst: "team1", //  Set to "team1" or "team2" to select who starts the veto. Any other values will default to team1 starting.
		SideType: "standard", // Either "standard", "always_knife", or "never_knife"

		// These values wrap mp_teamprediction_pct and mp_teamprediction_txt.
		// You can exclude these if you don't want those cvars set.
		// FavoredTeamPercentageText:"", //
		// FavoredTeamPercentageTeam1 : 50, //

		Maplist: m.VetoMapPoolJSON,
		// MapSides: "" // ??
	}
	//cfg.Spectators = make(map[string]string)
	//cfg.Spectators["STEAM_1:1:....."] = ""

	cfg.Team1.Flag = team1.Flag
	//cfg.Team1.Logo = ""
	cfg.Team1.Name = team1.Name
	cfg.Team1.Tag = team1.Tag
	// Any of the 3 formats (steam2, steam3, steam64 profile) are acceptable.
	// Note: the "players" section may be skipped if you set get5_check_auths to 0,
	// but that is not recommended. You can also set player names that will be forced here.
	// If you don't want to force player names, just use an empty quote "".
	cfg.Team1.Players = team1.SteamIDsJSON
	// cfg.Team1.Players = make(map[string]string)
	// cfg.Team1.Players["STEAM_0:1:52245092"] = "splewis"

	cfg.Team2.Flag = team2.Flag
	//cfg.Team2.Logo = ""
	cfg.Team2.Name = team2.Name
	cfg.Team2.Tag = team2.Tag
	cfg.Team2.Players = team2.SteamIDsJSON

	cfg.Cvars = make(map[string]string)
	commands := strings.Split(m.Cvars, "\n")
	for _, v := range commands {
		command := strings.Split(v, " ")
		if len(command) == 0 {
			// empty command...
		} else if len(command) != 2 {
			log.Printf("Something went wrong with match command. ERR : %v\n", command)
		} else {
			cfg.Cvars[command[0]] = command[1]
		}
	}
	cfg.Cvars["get5_web_api_url"] = fmt.Sprintf("http://%s/api/v1", config.Cnf.HOST)
	// cfg.Cvars["hostname"] = fmt.Sprintf("Match Server #1")

	return &cfg, nil
}

// GetMapStat Gets each map stat data as "MapStatsData" struct array.
func (m *MatchData) GetMapStat() (*MatchData, error) {
	SQLAccess.Gorm.Limit(7).Where("match_id = ?", m.ID).Find(&m.MapStats)
	return m, nil
}

// MapStatsData MapStatsData struct for map_stats table.
type MapStatsData struct {
	ID         int           `gorm:"primary_key" gorm:"column:id"`
	MatchID    int           `gorm:"column:match_id" gorm:"ForeignKey:match_id"`
	MapNumber  int           `gorm:"column:map_number"`
	MapName    string        `gorm:"column:map_name"`
	StartTime  sql.NullTime  `gorm:"column:start_time"`
	EndTime    sql.NullTime  `gorm:"column:end_time"`
	Winner     sql.NullInt32 `gorm:"column:winner"`
	Team1Score int           `gorm:"column:team1_score"`
	Team2Score int           `gorm:"column:team2_score"`

	User UserData `gorm:"ASSOCIATION_FOREIGNKEY:user_id"`
}

// TableName declairation for GORM
func (m *MapStatsData) TableName() string {
	return "map_stats"
}

// GetOrCreate Get or register mapstats data.
func (m *MapStatsData) GetOrCreate(matchID int, MapNumber int, mapname string) (*MapStatsData, error) {
	Match := MatchData{}
	MatchRecord := SQLAccess.Gorm.First(&Match, matchID)
	if MatchRecord.RecordNotFound() {
		return nil, fmt.Errorf("Match not found")
	}
	if MapNumber >= Match.MaxMaps {
		return nil, fmt.Errorf("MapNumber is greater than max map number")
	}
	m.MatchID = matchID
	m.MapNumber = MapNumber
	m.MapName = mapname
	MapStatsRecord := SQLAccess.Gorm.Where("match_id = ? AND map_number = ?", matchID, MapNumber).First(&m)
	if MapStatsRecord.RecordNotFound() {
		m.MatchID = matchID
		m.MapNumber = MapNumber
		m.MapName = mapname
		m.StartTime.Scan(time.Now())
		m.Team1Score = 0
		m.Team2Score = 0
		SQLAccess.Gorm.Create(&m)
		log.Printf("Created MapStatsData : %v\n", m)
		return m, nil
	}
	return m, nil
}

// PlayerStatsData Player stats data struct for player_stats table.
type PlayerStatsData struct {
	ID               int    `gorm:"primary_key;column:id"`
	MatchID          int    `gorm:"column:match_id"`
	MapID            int    `gorm:"column:map_id"`
	TeamID           int    `gorm:"column:team_id"`
	SteamID          string `gorm:"column:steam_id;unique"`
	Name             string `gorm:"column:name"`
	Kills            int    `gorm:"column:kills"`
	Deaths           int    `gorm:"column:deaths"`
	Roundsplayed     int    `gorm:"column:roundsplayed"`
	Assists          int    `gorm:"column:assists"`
	FlashbangAssists int    `gorm:"column:flashbang_assists"`
	Teamkills        int    `gorm:"column:teamkills"`
	Suicides         int    `gorm:"column:suicides"`
	HeadshotKills    int    `gorm:"column:headshot_kills"`
	Damage           int    `gorm:"column:damage"`
	BombPlants       int    `gorm:"column:bomb_plants"`
	BombDefuses      int    `gorm:"column:bomb_defuses"`
	V1               int    `gorm:"column:v1"`
	V2               int    `gorm:"column:v2"`
	V3               int    `gorm:"column:v3"`
	V4               int    `gorm:"column:v4"`
	V5               int    `gorm:"column:v5"`
	K1               int    `gorm:"column:k1"`
	K2               int    `gorm:"column:k2"`
	K3               int    `gorm:"column:k3"`
	K4               int    `gorm:"column:k4"`
	K5               int    `gorm:"column:k5"`
	FirstdeathCT     int    `gorm:"column:firstdeath_Ct"`
	FirstdeathT      int    `gorm:"column:firstdeath_t"`
	FirstkillCT      int    `gorm:"column:firstkill_ct"`
	FirstkillT       int    `gorm:"column:firstkill_t"`

	User UserData `gorm:"ASSOCIATION_FOREIGNKEY:user_id"`
}

// TableName declairation for GORM
func (p *PlayerStatsData) TableName() string {
	return "player_stats"
}

// GetOrCreate Get or register player stats data into DB.
func (p *PlayerStatsData) GetOrCreate(matchID int, MapNumber int, steamid string) (*PlayerStatsData, error) {
	MapStats := &MapStatsData{}
	MapStatsRecord := SQLAccess.Gorm.Where("match_id = ? AND map_number = ?", matchID, MapNumber).First(&MapStats)
	if MapStatsRecord.RecordNotFound() {
		return nil, fmt.Errorf("MapStats not found")
	}
	// original get5-web restricts player per map stats length https://github.com/splewis/get5-web/blob/8c1012c9563583353b9486a6590227855547e275/get5/models.py#L558
	PlayerStats := &PlayerStatsData{}
	PlayerStatsRecord := SQLAccess.Gorm.Where("steam_id = ? AND match_id = ? AND map_id = ?", steamid, matchID, MapStats.ID).First(&PlayerStats)
	if PlayerStatsRecord.RecordNotFound() {
		PlayerStats.MatchID = matchID
		// PlayerStats.map_number  = MapStats.ID // Not exist..?? https://github.com/splewis/get5-web/blob/8c1012c9563583353b9486a6590227855547e275/get5/models.py#L566
		PlayerStats.SteamID = steamid
		PlayerStats.MapID = MapStats.ID
		SQLAccess.Gorm.Create(&PlayerStats)
	}
	return PlayerStats, nil
}

// GetSteamURL get player's Steam community URL by their steamid64.
func (p *PlayerStatsData) GetSteamURL() string {
	return fmt.Sprintf("http://steamcommunity.com/profiles/%s", p.SteamID)
}

// GetRating Get player's rating. Average datas are static tho.
func (p *PlayerStatsData) GetRating() float64 { // Rating value can be more accurate??
	var AverageKPR float64 = 0.679
	var AverageSPR float64 = 0.317
	var AverageRMK float64 = 1.277
	var KillRating float64 = float64(p.Kills) / float64(p.Roundsplayed) / AverageKPR
	var SurvivalRating float64 = float64(p.Roundsplayed-p.Deaths) / float64(p.Roundsplayed) / float64(AverageSPR)
	var killcount float64 = float64(p.K1 + 4*p.K2 + 9*p.K3 + 16*p.K4 + 25*p.K5)
	var RoundsWithMultipleKillsRating float64 = killcount / float64(p.Roundsplayed) / float64(AverageRMK)
	var rating float64 = (KillRating + 0.7*SurvivalRating + RoundsWithMultipleKillsRating) / 2.7
	return rating
}

// GetKDR Returns player's KDR(Kill/Deaths Ratio).
func (p *PlayerStatsData) GetKDR() float64 {
	if p.Deaths == 0 {
		return float64(p.Kills)
	}
	return float64(p.Kills) / float64(p.Deaths)
}

// GetHSP Returns player's HSP(HeadShot Percentage).
func (p *PlayerStatsData) GetHSP() float64 {
	if p.Kills == 0 {
		return math.NaN()
	}
	return float64(p.HeadshotKills) / float64(p.Kills) * 100
}

// GetADR Returns player's ADR(Average Damage per Round).
func (p *PlayerStatsData) GetADR() float64 {
	if p.Roundsplayed == 0 {
		return 0.0
	}
	return float64(p.Damage) / float64(p.Roundsplayed)
}

// GetFPR Returns player's FPR(Frags Per Round).
func (p *PlayerStatsData) GetFPR() float64 {
	if p.Roundsplayed == 0 {
		return 0.0
	}
	return float64(p.Kills) / float64(p.Roundsplayed)
}

// GetSteamName Get steam profile name by steamid64 via Steam web API
func GetSteamName(steamid uint64) (string, error) {
	summary, err := steam.GetPlayerSummaries(SteamAPIKey, steam.SteamID64(steamid))
	if err != nil {
		return "", err
	}
	return summary.DisplayName, nil
}

// MetricsData Struct metrics analysys.
type MetricsData struct {
	RegisteredUsers    int `json:"users"`
	SavedTeams         int `json:"saved_teams"`
	MatchesCreated     int `json:"matches_created"`
	CompletedMatches   int `json:"completed_matches"`
	ServersAdded       int `json:"servers_added"`
	MapsWithStatsSaved int `json:"maps_with_stats"`
	UniquePlayers      int `json:"unique_players"`
}

// GetMetrics Get Each table's count.
func GetMetrics() MetricsData {
	var result MetricsData

	SQLAccess.Gorm.Table("user").Count(&result.RegisteredUsers)
	SQLAccess.Gorm.Table("team").Count(&result.SavedTeams)
	SQLAccess.Gorm.Table("match").Count(&result.MatchesCreated)
	SQLAccess.Gorm.Table("match").Where("end_time IS NOT NULL").Count(&result.CompletedMatches)
	SQLAccess.Gorm.Table("game_server").Count(&result.ServersAdded)
	SQLAccess.Gorm.Table("map_stats").Count(&result.MapsWithStatsSaved)
	SQLAccess.Gorm.Table("player_stats").Count(&result.UniquePlayers)

	return result
}

func checkIP(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		log.Printf("%v is not an IPv4 address\n", ip)
		return false
	}
	return true
}

// TODO ADD JSON Meta data...

// RoundStatsData RoundStatsData struct for map_stats table.
type RoundStatsData struct {
	ID                   int            `gorm:"primary_key" gorm:"column:id"`
	MatchID              int            `gorm:"column:match_id"`
	MapNumber            int            `gorm:"column:map_number"`
	Winner               sql.NullString `gorm:"column:winner"`
	WinnerSide           sql.NullString `gorm:"column:winner_side"`
	FirstKillerSteamID   sql.NullString `gorm:"column:first_killer_steamid"`
	FirstVictimSteamID   sql.NullString `gorm:"column:first_victim_steamid"`
	SecondKillerSteamID  sql.NullString `gorm:"column:second_killer_steamid"`
	SecondVictimSteamID  sql.NullString `gorm:"column:second_victim_steamid"`
	ThirdKillerSteamID   sql.NullString `gorm:"column:third_killer_steamid"`
	ThirdVictimSteamID   sql.NullString `gorm:"column:third_victim_steamid"`
	FourthKillerSteamID  sql.NullString `gorm:"column:fourth_killer_steamid"`
	FourthVictimSteamID  sql.NullString `gorm:"column:fourth_victim_steamid"`
	FifthKillerSteamID   sql.NullString `gorm:"column:fifth_killer_steamid"`
	FifthVictimSteamID   sql.NullString `gorm:"column:fifth_victim_steamid"`
	SixthKillerSteamID   sql.NullString `gorm:"column:sixth_killer_steamid"`
	SixthVictimSteamID   sql.NullString `gorm:"column:sixth_victim_steamid"`
	SeventhKillerSteamID sql.NullString `gorm:"column:seventh_killer_steamid"`
	SeventhVictimSteamID sql.NullString `gorm:"column:seventh_victim_steamid"`
	EighthKillerSteamID  sql.NullString `gorm:"column:eighth_killer_steamid"`
	EighthVictimSteamID  sql.NullString `gorm:"column:eighth_victim_steamid"`
	NinthKillerSteamID   sql.NullString `gorm:"column:ninth_killer_steamid"`
	NinthVictimSteamID   sql.NullString `gorm:"column:ninth_victim_steamid"`
	TenthKillerSteamID   sql.NullString `gorm:"column:tenth_killer_steamid"`
	TenthVictimSteamID   sql.NullString `gorm:"column:tenth_victim_steamid"`

	Match MatchData    `gorm:"ASSOCIATION_FOREIGNKEY:match_id"`
	Map   MapStatsData `gorm:"ASSOCIATION_FOREIGNKEY:map_id"`
}

// TableName declairation for GORM
func (r *RoundStatsData) TableName() string {
	return "round_stats"
}

// GetOrCreate Get or register round stats data into DB.
func (r *RoundStatsData) GetOrCreate(matchID int, MapNumber int) (*RoundStatsData, error) {
	log.Printf("Registering round info : %v\n", *r)
	rec := SQLAccess.Gorm.FirstOrCreate(r)
	log.Printf("rec : %v\n", rec)
	if len(rec.GetErrors()) > 1 {
		return nil, rec.GetErrors()[0]
	}
	return r, nil
}

// Register round stats data into DB.
func (r *RoundStatsData) Register(matchID int, MapNumber int) (*RoundStatsData, error) {
	match := &MatchData{}
	match, err := match.Get(matchID)
	if err != nil {
		log.Printf("Failed to get match : %v\n", err)
		return nil, err
	}
	log.Printf("Match : %v\n", match)
	if !match.Live() {
		log.Println("Match is not Live")
		return nil, fmt.Errorf("Match is not Live")
	}
	log.Printf("Registering round info : %v\n", *r)
	rec := SQLAccess.Gorm.Create(r)
	log.Printf("rec : %v\n", rec)
	if len(rec.GetErrors()) > 1 {
		return nil, rec.GetErrors()[0]
	}
	return r, nil
}
