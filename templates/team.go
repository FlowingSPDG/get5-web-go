// This file is generated by gorazor 2.0.1
// DON'T modified manually
// Should edit source file and re-generate: templates/team.gohtml

package templates

import (
	"github.com/FlowingSPDG/get5-web-go/src/db"
	"github.com/FlowingSPDG/get5-web-go/templates/layout"
	"github.com/sipin/gorazor/gorazor"
	"io"
	"strings"
)

// Team generates templates/team.gohtml
func Team(u *db.TeamPageData) string {
	var _b strings.Builder
	RenderTeam(&_b, u)
	return _b.String()
}

// RenderTeam render templates/team.gohtml
func RenderTeam(_buffer io.StringWriter, u *db.TeamPageData) {

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

		_buffer.WriteString("<div id=\"content\">\n\n  <div class=\"container\">\n    <h1>\n      ")
		_buffer.WriteString((u.Team.GetFlagHTML(1.0)))
		_buffer.WriteString(" ")
		_buffer.WriteString((u.Team.Name))
		_buffer.WriteString(" ")
		_buffer.WriteString((u.Team.GetLogoHtml(1.0)))
		_buffer.WriteString("\n      ")
		if u.Team.CanEdit(u.User.ID) {

			_buffer.WriteString("<div class=\"pull-right\">\n        <a href=\"/team/")
			_buffer.WriteString(gorazor.HTMLEscape(u.Team.ID))
			_buffer.WriteString("/edit\" class=\"btn btn-primary btn-xs\">Edit</a>\n      </div>")

		}
		_buffer.WriteString("\n    </h1>\n\n    <br>\n\n    <div class=\"panel panel-default\">\n      <div class=\"panel-heading\">Players</div>\n      <div class=\"panel-body\">\n        ")

		players, _ := u.Team.GetPlayers()

		_buffer.WriteString("\n        ")
		for i := 0; i < len(u.Team.Auths); i++ {
			p := players[i]

			_buffer.WriteString("<a href=\"http://steamcommunity.com/profiles/")
			_buffer.WriteString(gorazor.HTMLEscape(u.Team.Auths[i]))
			_buffer.WriteString("\" class=\"col-sm-offset-0\"> ")
			_buffer.WriteString(gorazor.HTMLEscape(u.Team.Auths[i]))
			_buffer.WriteString(" </a>")
			_buffer.WriteString(gorazor.HTMLEscape(p.Name))

			_buffer.WriteString("<br>\n        ")
		}
		_buffer.WriteString("\n      </div>\n    </div>\n\n\n    <div class=\"panel panel-default\">\n      <div class=\"panel-heading\">Recent Matches</div>\n        <div class=\"panel-body\">\n          ")

		matches := u.Team.GetRecentMatches(100)

		_buffer.WriteString("\n          ")
		for i := 0; i < len(matches); i++ {
			m := matches[i]
			matchresult, _ := u.Team.GetVSMatchResult(int(m.ID))

			_buffer.WriteString("<a href=\"/match/")
			_buffer.WriteString(gorazor.HTMLEscape(m.ID))
			_buffer.WriteString("\">#")
			_buffer.WriteString(gorazor.HTMLEscape(m.ID))
			_buffer.WriteString("</a>")
			_buffer.WriteString((":"))
			_buffer.WriteString(gorazor.HTMLEscape(matchresult))

			_buffer.WriteString("<br>\n          ")
		}
		_buffer.WriteString("\n          \n      </div>\n    </div>\n\n  </div>\n  <br>\n\n</div>")

	}

	layout.RenderBase(_buffer, _body, _menu, _content)
}
