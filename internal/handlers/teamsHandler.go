package handlers

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func (hm *HandlersManager) TeamsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	roasters, err := hm.servicesManager.GetTeamsRoastersList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hm.logger.Log(err)
		return
	}
	data["TeamsRoasters"] = roasters
	net.RenderTemplate(w, r, "teams.html", data)
}
