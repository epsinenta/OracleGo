package handlers

import (
	"OracleGo/internal/entities"
	_ "fmt"
	"net/http"
)

func (hm *HandlersManager) StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	heroes, err := hm.servicesManager.Statistics.GetAllHeroesWinrates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hm.logger.Log(err)
		return
	}
	data["HeroesWinrates"] = heroes

	isLoggedIn, err := getBoolFromContext(r, entities.AuthStatusContextKey{})
	if err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data[entities.IsLoggedInKey] = isLoggedIn

	if err := hm.templates.ExecuteTemplate(w, "statistics.html", data); err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
