package ml

import (
	"OracleGo/internal/net"
	"OracleGo/internal/statistics"
	_ "fmt"
	"net/http"
)

func PredictionHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	dbTeamManager, _ := statistics.NewTeamsDatabaseManager()
	teams, _ := dbTeamManager.GetTeamsList()
	data["Team"] = teams

	dbHeroManager, _ := statistics.NewHeroesDatabaseManager()
	heroes, _ := dbHeroManager.GetHeroesNameList()
	data["Hero"] = heroes

	net.RenderTemplate(w, r, "prediction.html", data)
}
