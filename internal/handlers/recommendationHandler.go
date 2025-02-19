package handlers

import (
	"OracleGo/internal/entities"
	_ "fmt"
	"net/http"
)

func (hm *HandlersManager) RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	isLoggedIn, err := getBoolFromContext(r, entities.AuthStatusContextKey{})
	if err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data[entities.IsLoggedInKey] = isLoggedIn

	if err := hm.templates.ExecuteTemplate(w, "prediction.html", data); err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
