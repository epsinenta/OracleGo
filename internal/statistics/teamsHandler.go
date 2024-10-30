package statistics

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	dbManager, _ := NewTeamsDatabaseManager()
	roasters, _ := dbManager.GetTeamsRoastersList()
	data["TeamsRoaters"] = roasters
	net.RenderTemplate(w, r, "teams.html", data)
}
