// This file is generated by gorazor 2.0.1
// DON'T modified manually
// Should edit source file and re-generate: templates/match.gohtml

package templates

import (
	db "github.com/FlowingSPDG/get5-web-go/src/db"
	"github.com/FlowingSPDG/get5-web-go/templates/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strconv"
	"strings"
)

// Match generates templates/match.gohtml
func Match(u *db.MatchPageData) string {
	var _b strings.Builder
	RenderMatch(&_b, u)
	return _b.String()
}

// RenderMatch render templates/match.gohtml
func RenderMatch(_buffer io.StringWriter, u *db.MatchPageData) {

	_body := func(_buffer io.StringWriter) {

	}

	_menu := func(_buffer io.StringWriter) {

		if u.LoggedIn == true {

			_buffer.WriteString("<li><a id=\"mymatches\" href=\"/mymatches\">My Matches</a></li>")

			_buffer.WriteString("<li><a id=\"match_create\" href=\"/match/create\">Create a Match</a></li>")

			_buffer.WriteString("<li><a id=\"myteams\" href=\"/myteams\">My Teams</a></li>")

			_buffer.WriteString("<li><a id=\"team_create\" href=\"/team/create\">Create a Team</a></li>")

			_buffer.WriteString("<li><a id=\"myservers\" href=\"/myservers\">My Servers</a></li>")

			_buffer.WriteString("<li><a id=\"server_create\" href=\"/server/create\">Add a Server</a></li>")

			_buffer.WriteString("<li><a href=\"/logout\">Logout</a></li>")

		} else {

			_buffer.WriteString("<li><a href=\"/login\"><img src=\"/static/img/login_small.png\" height=\"18\"/></a></li>")

		}

	}

	_content := func(_buffer io.StringWriter) {

		_buffer.WriteString("<div id=\"content\">\n  <div class=\"container\">\n    <h1>\n      ")

		team1, _ := u.Match.GetTeam1()
		team2, _ := u.Match.GetTeam2()
		mapstats, _ := u.Match.GetMapStat()
		var timeformat = "2006/01/02 15:04"

		_buffer.WriteString("\n      ")
		_buffer.WriteString((team1.GetLogoOrFlagHtml(1.0, team2)))
		_buffer.WriteString(" <a href=\"/team/")
		_buffer.WriteString(gorazor.HTMLEscape(team1.ID))
		_buffer.WriteString("\"> ")
		_buffer.WriteString(gorazor.HTMLEscape(team1.Name))
		_buffer.WriteString("</a>\n      ")
		_buffer.WriteString(gorazor.HTMLEscape(u.Match.Team1Score))
		_buffer.WriteString("\n      ")

		if u.Match.Team1Score < u.Match.Team2Score {

			_buffer.WriteString(("<"))

		} else if u.Match.Team1Score == u.Match.Team2Score {

			_buffer.WriteString(("=="))

		} else {

			_buffer.WriteString((">"))

		}

		_buffer.WriteString("\n      ")
		_buffer.WriteString(gorazor.HTMLEscape(u.Match.Team2Score))
		_buffer.WriteString("\n      ")
		_buffer.WriteString((team2.GetLogoOrFlagHtml(1.0, team1)))
		_buffer.WriteString(" <a href=\"/team/")
		_buffer.WriteString(gorazor.HTMLEscape(team2.ID))
		_buffer.WriteString("\"> ")
		_buffer.WriteString(gorazor.HTMLEscape(team2.Name))
		_buffer.WriteString("</a>\n\n      ")
		if u.AdminAccess && u.Match.Live() || u.Match.Pending() {

			_buffer.WriteString("<div class=\"dropdown dropdown-header pull-right\">\n        <button class=\"btn btn-default dropdown-toggle\" type=\"button\" id=\"dropdownMenu1\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"true\">\n          Admin tools\n          <span class=\"caret\"></span>\n        </button>\n        <ul class=\"dropdown-menu\" aria-labelledby=\"dropdownMenu1\">\n          if u.Match.Live(){\n            <li><a id=\"pause\" href=\"{{request.path}}/pause\">Pause match</a></li>\n            <li><a id=\"unpause\" href=\"{{request.path}}/unpause\">Unpause match</a></li>\n          }\n          <li><a id=\"addplayer_team1\" href=\"#\">Add player to team1</a></li>\n          <li><a id=\"addplayer_team2\" href=\"#\">Add player to team2</a></li>\n          <li><a id=\"addplayer_spec\" href=\"#\">Add player to specator list</a></li>\n          <li><a id=\"rcon_command\" href=\"#\">Send rcon command</a></li>\n          <li role=\"separator\" class=\"divider\"></li>\n          <li><a id=\"backup_manager\" href=\"{{request.path}}/backup\">Load a backup file</a></li>\n          <li><a href=\"{{request.path}}/cancel\">Cancel match</a></li>\n        </ul>\n      </div>")

		}
		_buffer.WriteString("\n    </h1>\n\n\n    <br>\n    ")
		if u.Match.Cancelled {

			_buffer.WriteString("<div class=\"alert alert-danger\" role=\"alert\">\n      <span class=\"glyphicon glyphicon-exclamation-sign\" aria-hidden=\"true\"></span>\n      <span class=\"sr-only\">Error:</span>\n        This match has been cancelled.\n      </div>")

		}
		_buffer.WriteString("\n\n    ")
		if u.Match.Forfeit {
			loser, _ := u.Match.GetLoser()

			_buffer.WriteString("<div class=\"alert alert-warning\" role=\"alert\">\n      <span class=\"glyphicon glyphicon-exclamation-sign\" aria-hidden=\"true\"></span>\n      <span class=\"sr-only\">Error:</span>\n        This match was forfeit by ")
			_buffer.WriteString(gorazor.HTMLEscape(loser))
			_buffer.WriteString(" .\n      </div>")

		}
		_buffer.WriteString("\n\n\n    ")
		if u.Match.StartTime.Valid {

			starttime := u.Match.StartTime.Time.Format(timeformat)

			_buffer.WriteString("<p>Started at ")
			_buffer.WriteString(gorazor.HTMLEscape(starttime))
			_buffer.WriteString("</p>")

		} else {

			_buffer.WriteString("<div class=\"panel panel-default\" role=\"alert\">\n      <div class=\"panel-body\">\n        This match is pending start.\n      </div>\n    </div>\n    ")
		}
		_buffer.WriteString("\n\n\n    ")
		if u.Match.EndTime.Valid {

			endtime := u.Match.EndTime.Time.Format(timeformat)

			_buffer.WriteString("<p>Ended at ")
			_buffer.WriteString(gorazor.HTMLEscape(endtime))
			_buffer.WriteString("</p>")

		}
		_buffer.WriteString("\n\n    ")
		for i := 0; i < len(mapstats); i++ {

			mTeam1Score := strconv.Itoa(mapstats[i].Team1Score)
			mTeam2Score := strconv.Itoa(mapstats[i].Team2Score)

			_buffer.WriteString("<br>\n    <div class=\"panel panel-primary\">\n      <div class=\"panel-heading\">\n        Map ")
			_buffer.WriteString(gorazor.HTMLEscape(mapstats[i]))
			_buffer.WriteString(".MapNumber+1: ")
			_buffer.WriteString(gorazor.HTMLEscape(mapstats[i]))
			_buffer.WriteString(".MapName,\n        ")
			_buffer.WriteString(gorazor.HTMLEscape(team1.Name))
			_buffer.WriteString(" \n        ")

			if u.Match.Team1Score < u.Match.Team2Score {

				_buffer.WriteString(("<"))

			} else if u.Match.Team1Score == u.Match.Team2Score {

				_buffer.WriteString(("=="))

			} else {

				_buffer.WriteString((">"))

			}

			_buffer.WriteString("\n        ")
			_buffer.WriteString(gorazor.HTMLEscape(team2.Name))
			_buffer.WriteString(",\n        ")
			_buffer.WriteString(gorazor.HTMLEscape(mTeam1Score))
			_buffer.WriteString((":"))
			_buffer.WriteString(gorazor.HTMLEscape(mTeam2Score))
			_buffer.WriteString("\n      </div>\n\n      <div class=\"panel-body\">\n        <p>Started at ")
			_buffer.WriteString(gorazor.HTMLEscape(mapstats[i]))
			_buffer.WriteString(".StartTime</p>\n\n        ")
			if mapstats[i].EndTime.Valid {

				_buffer.WriteString("<p>Ended at ")
				_buffer.WriteString(gorazor.HTMLEscape(mapstats[i]))
				_buffer.WriteString(".EndTime</p>")

			}
			_buffer.WriteString("\n\n        <table class=\"table table-hover\">\n          <thead>\n            <tr>\n              <th>Player</th>\n              <th class=\"text-center\">Kills</th>\n              <th class=\"text-center\">Deaths</th>\n              <th class=\"text-center\">Assists</th>\n              <th class=\"text-center\">Flash assists</th>\n              <th class=\"text-center\">1v1</th>\n              <th class=\"text-center\">1v2</th>\n              <th class=\"text-center\">1v3</th>\n              <th class=\"text-center\">Rating</th>\n              <th class=\"text-center\"><acronym title=\"Frags per round\">FPR</acronym></th>\n              <th class=\"text-center\"><acronym title=\"Average damage per round\">ADR</acronym></th>\n              <th class=\"text-center\"><acronym title=\"Headshot percentage\">HSP</acronym></th>\n            </tr>\n          </thead>\n          <tbody>\n          ")
			players1, _ := team1.GetPlayers()
			_buffer.WriteString("\n            ")
			for i := 0; i < len(players1); i++ {
				player := players1[i]
				db.SQLAccess.Gorm.Where("team_id = ?", team1.ID).Find(&player)
				if player.Roundsplayed > 0 {

					_buffer.WriteString("<tr>\n                <td> <a href=\"")
					_buffer.WriteString(gorazor.HTMLEscape(player.GetSteamURL()))
					_buffer.WriteString("\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Name))
					_buffer.WriteString(" </a></td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Kills))
					_buffer.WriteString("  </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Deaths))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Assists))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Flashbang_assists))
					_buffer.WriteString(" </td>\n\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V1))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V2))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V3))
					_buffer.WriteString(" </td>\n\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(strconv.Itoa(int(player.GetRating()))))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> {{ player.get_fpr() | round(2) }} </td>\n                <td class=\"text-center\"> {{ player.get_adr() | round(1) }} </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(strconv.Itoa(int(player.GetHSP()))))
					_buffer.WriteString(" </td>\n              </tr>")

				}
			}
			_buffer.WriteString("\n          ")
			players2, _ := team2.GetPlayers()
			_buffer.WriteString("\n            ")
			for i := 0; i < len(players2); i++ {
				player := players2[i]
				db.SQLAccess.Gorm.Where("team_id = ?", team2.ID).Find(&player)
				if player.Roundsplayed > 0 {

					_buffer.WriteString("<tr>\n                <td> <a href=\"")
					_buffer.WriteString(gorazor.HTMLEscape(player.GetSteamURL()))
					_buffer.WriteString("\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Name))
					_buffer.WriteString(" </a></td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Kills))
					_buffer.WriteString("  </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Deaths))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Assists))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.Flashbang_assists))
					_buffer.WriteString(" </td>\n\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V1))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V2))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(player.V3))
					_buffer.WriteString(" </td>\n\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(strconv.Itoa(int(player.GetRating()))))
					_buffer.WriteString(" </td>\n                <td class=\"text-center\"> {{ player.get_fpr() | round(2) }} </td>\n                <td class=\"text-center\"> {{ player.get_adr() | round(1) }} </td>\n                <td class=\"text-center\"> ")
					_buffer.WriteString(gorazor.HTMLEscape(strconv.Itoa(int(player.GetHSP()))))
					_buffer.WriteString(" </td>\n              </tr>")

				}
			}
			_buffer.WriteString("\n          </tbody>\n        </table>\n      </div>\n    ")
		}
		_buffer.WriteString("\n    </div>\n  </div>\n\n\n  <br>\n</div>\n\n<script>\n\njQuery(\"#addplayer_team1\").click(function(e) {\n    var input = prompt(\"Please enter a steamid to add to {{team1.name}}\", \"\");\n    if (input != null) {\n      window.location.href = \"{{request.path}}/adduser?team=team1&auth=\" + encodeURIComponent(input);\n    }\n});\n\njQuery(\"#addplayer_team2\").click(function(e) {\n    var input = prompt(\"Please enter a steamid to add to {{team2.name}}\", \"\");\n    if (input != null) {\n      window.location.href = \"{{request.path}}/adduser?team=team2&auth=\" + encodeURIComponent(input);\n    }\n});\n\njQuery(\"#addplayer_spec\").click(function(e) {\n    var input = prompt(\"Please enter a steamid to add to the spectators list\", \"\");\n    if (input != null) {\n      window.location.href = \"{{request.path}}/adduser?team=spec&auth=\" + encodeURIComponent(input);\n    }\n});\n\njQuery(\"#rcon_command\").click(function(e) {\n    var input = prompt(\"Enter a command to send\", \"\");\n    if (input != null) {\n      window.location.href = \"{{request.path}}/rcon?command=\" + encodeURIComponent(input);\n    }\n});\n</script>")
	}

	layout.RenderBase(_buffer, _body, _menu, _content)
}
