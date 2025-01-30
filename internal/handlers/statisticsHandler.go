package handlers

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func (hm *HandlersManager) StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	heroes, err := hm.servicesManager.GetAllHeroesWinrates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hm.logger.Log(err)
		return
	}
	data["HeroesWinrates"] = heroes
	net.RenderTemplate(w, r, "statistics.html", data)
}
