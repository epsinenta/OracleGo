package statistics

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	dbManager, _ := NewHeroesDatabaseManager()
	heroes, _ := dbManager.GetAllHeroesWinrates()
	data["HeroesWinrates"] = heroes
	net.RenderTemplate(w, r, "statistics.html", data)
}
