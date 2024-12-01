package ml

import (
	"OracleGo/internal/net"
	"OracleGo/internal/statistics"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	dbTeamManager   *statistics.TeamsDatabaseManager
	dbHeroManager   *statistics.HeroesDatabaseManager
	dbPlayerManager *statistics.PlayersDatabaseManager
)

func init() {
	var err error
	// Инициализация менеджеров баз данных один раз
	dbTeamManager, err = statistics.NewTeamsDatabaseManager()
	if err != nil {
		fmt.Printf("Failed to create teams database manager: %v\n", err)
		return
	}
	dbHeroManager, err = statistics.NewHeroesDatabaseManager()
	if err != nil {
		fmt.Printf("Failed to create heroes database manager: %v\n", err)
		return
	}
	dbPlayerManager, err = statistics.NewPlayersDatabaseManager()
	if err != nil {
		fmt.Printf("Failed to create players database manager: %v\n", err)
		return
	}
}

func PredictionHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		team1 := r.FormValue("team1")
		team2 := r.FormValue("team2")

		var team1Players, team2Players []string
		teamRosters, err := dbTeamManager.GetTeamsRoastersList()
		if err != nil {
			http.Error(w, "Failed to get teams rosters list", http.StatusInternalServerError)
			return
		}

		team1Processed := false
		team2Processed := false

		for _, roster := range teamRosters {
			teamName := roster.Team.GetValue()

			if !team1Processed && teamName == team1 {
				for _, player := range roster.Players {
					team1Players = append(team1Players, player.GetValue())
				}
				team1Processed = true
			} else if !team2Processed && teamName == team2 {
				for _, player := range roster.Players {
					team2Players = append(team2Players, player.GetValue())
				}
				team2Processed = true
			}

			if team1Processed && team2Processed {
				break
			}
		}

		team1Heroes := []string{
			r.FormValue("team1-hero1"),
			r.FormValue("team1-hero2"),
			r.FormValue("team1-hero3"),
			r.FormValue("team1-hero4"),
			r.FormValue("team1-hero5"),
		}
		team2Heroes := []string{
			r.FormValue("team2-hero1"),
			r.FormValue("team2-hero2"),
			r.FormValue("team2-hero3"),
			r.FormValue("team2-hero4"),
			r.FormValue("team2-hero5"),
		}

		allHeroes := append(team1Heroes, team2Heroes...)
		var heroes []statistics.Hero
		for _, heroName := range allHeroes {
			heroes = append(heroes, statistics.Hero{Value: heroName})
		}
		heroWinrates, err := dbHeroManager.GetHeroesWinrates(heroes)
		if err != nil {
			http.Error(w, "Failed to get heroes winrates", http.StatusInternalServerError)
			return
		}

		var team1HeroObjs, team2HeroObjs []statistics.Hero
		for _, name := range team1Heroes {
			team1HeroObjs = append(team1HeroObjs, statistics.Hero{Value: name})
		}
		for _, name := range team2Heroes {
			team2HeroObjs = append(team2HeroObjs, statistics.Hero{Value: name})
		}
		heroCounterPicks, err := dbHeroManager.GetHeroesCounterPicks(team1HeroObjs, team2HeroObjs)
		if err != nil {
			http.Error(w, "Failed to get heroes counter picks", http.StatusInternalServerError)
			return
		}

		var players []statistics.Player
		for _, playerName := range append(team1Players, team2Players...) {
			players = append(players, statistics.Player{Value: playerName})
		}
		playerWinrates, err := dbPlayerManager.GetPlayerOnHeroWinrate(players, heroes)
		if err != nil {
			http.Error(w, "Failed to get player winrates on heroes", http.StatusInternalServerError)
			return
		}

		playerGameCounts, err := dbPlayerManager.GetPlayerCountOnHero(players, heroes)
		if err != nil {
			http.Error(w, "Failed to get player game counts on heroes", http.StatusInternalServerError)
			return
		}

		var rowData []string
		rowData = append(rowData, team1, team2)
		rowData = append(rowData, team1Players...)
		rowData = append(rowData, team2Players...)
		rowData = append(rowData, team1Heroes...)
		rowData = append(rowData, team2Heroes...)

		for _, winrate := range heroWinrates {
			rowData = append(rowData, fmt.Sprintf("%.2f", winrate.Winrate))
		}

		for _, team1Hero := range team1Heroes {
			for _, team2Hero := range team2Heroes {
				found := false
				for _, counter := range heroCounterPicks {
					if counter.FirstHero.Value == team1Hero && counter.SecondHero.Value == team2Hero {
						rowData = append(rowData, fmt.Sprintf("%.2f", counter.CounterPick))
						found = true
						break
					}
				}
				if !found {
					rowData = append(rowData, "0.00")
				}
			}
		}

		for i, player := range team1Players {
			hero := team1Heroes[i]
			found := false
			for _, winrate := range playerWinrates {
				if winrate.Player.Value == player && winrate.Hero.Value == hero {
					rowData = append(rowData, fmt.Sprintf("%.2f", winrate.Winrate))
					found = true
					break
				}
			}
			if !found {
				rowData = append(rowData, "0.00")
			}
		}

		for i, player := range team2Players {
			hero := team2Heroes[i]
			found := false
			for _, winrate := range playerWinrates {
				if winrate.Player.Value == player && winrate.Hero.Value == hero {
					rowData = append(rowData, fmt.Sprintf("%.2f", winrate.Winrate))
					found = true
					break
				}
			}
			if !found {
				rowData = append(rowData, "0.00")
			}
		}

		for i, player := range team1Players {
			hero := team1Heroes[i]
			found := false
			for _, gameCount := range playerGameCounts {
				if gameCount.Player.Value == player && gameCount.Hero.Value == hero {
					rowData = append(rowData, fmt.Sprintf("%d", gameCount.Count))
					found = true
					break
				}
			}
			if !found {
				rowData = append(rowData, "0")
			}
		}

		for i, player := range team2Players {
			hero := team2Heroes[i]
			found := false
			for _, gameCount := range playerGameCounts {
				if gameCount.Player.Value == player && gameCount.Hero.Value == hero {
					rowData = append(rowData, fmt.Sprintf("%d", gameCount.Count))
					found = true
					break
				}
			}
			if !found {
				rowData = append(rowData, "0")
			}
		}

		outputRow := strings.Join(rowData, ",")

		if err := os.WriteFile("internal/ml/scripts/row.txt", []byte(outputRow), 0644); err != nil {
			http.Error(w, "Failed to write to row.txt", http.StatusInternalServerError)
			return
		}

		cmd := exec.Command("python", "internal/ml/scripts/run_model.py")
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "Failed to run model script", http.StatusInternalServerError)
			fmt.Printf("Script error: %s\n", err)
			fmt.Printf("Script output: %s\n", output)
			return
		}

		file2, err := os.Open("internal/ml/scripts/prediction_result.json")
		if err != nil {
			http.Error(w, "Failed to open prediction result file", http.StatusInternalServerError)
			return
		}
		defer file2.Close()

		var resultData map[string]interface{}
		decoder := json.NewDecoder(file2)
		if err := decoder.Decode(&resultData); err != nil {
			http.Error(w, "Failed to parse prediction result", http.StatusInternalServerError)
			return
		}

		data["PredictionResult"] = resultData["prediction"]
		data["PredictionProbability"] = resultData["probability"]
	}

	if teams, err := dbTeamManager.GetTeamsList(); err == nil {
		data["Team"] = teams
	}
	if heroList, err := dbHeroManager.GetHeroesNameList(); err == nil {
		data["Hero"] = heroList
	}

	net.RenderTemplate(w, r, "prediction.html", data)
}
