package handlers

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/net"
	"OracleGo/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func (hm *HandlersManager) PredictionHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		team1 := r.FormValue("team1")
		team2 := r.FormValue("team2")

		var team1Players, team2Players []string
		teamRosters, err := hm.servicesManager.GetTeamsRoastersList()
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
		var heroes []entities.Hero
		for _, heroName := range allHeroes {
			heroes = append(heroes, entities.Hero{Value: heroName})
		}
		heroWinrates, err := hm.servicesManager.GetHeroesWinrates(heroes)
		if err != nil {
			http.Error(w, "Failed to get heroes winrates", http.StatusInternalServerError)
			return
		}

		var team1HeroObjs, team2HeroObjs []entities.Hero
		for _, name := range team1Heroes {
			team1HeroObjs = append(team1HeroObjs, entities.Hero{Value: name})
		}
		for _, name := range team2Heroes {
			team2HeroObjs = append(team2HeroObjs, entities.Hero{Value: name})
		}
		heroCounterPicks, err := hm.servicesManager.GetHeroesCounterPicks(team1HeroObjs, team2HeroObjs)
		if err != nil {
			http.Error(w, "Failed to get heroes counter picks", http.StatusInternalServerError)
			return
		}

		var players []entities.Player
		for _, playerName := range append(team1Players, team2Players...) {
			players = append(players, entities.Player{Value: playerName})
		}

		//playerWinrates, err := hm.servicesManager.GetPlayersWinrate(players)
		// if err != nil {
		// 	http.Error(w, "Failed to get player winrates on heroes", http.StatusInternalServerError)
		// 	return
		// }

		playerOnHeroWinrates, err := hm.servicesManager.GetPlayerOnHeroWinrate(players, heroes)
		if err != nil {
			http.Error(w, "Failed to get player winrates on heroes", http.StatusInternalServerError)
			return
		}

		playerGameCounts, err := hm.servicesManager.GetPlayerCountOnHero(players, heroes)
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

		////////////////////////////

		// for _, winrate := range playerWinrates {
		//	rowData = append(rowData, fmt.Sprintf("%.2f", winrate.Winrate))
		//}

		for i := 0; i < 10; i++ {
			rowData = append(rowData, "0 ")
		}

		////////////////////////////

		for i, player := range team1Players {
			hero := team1Heroes[i]
			found := false
			for _, winrate := range playerOnHeroWinrates {
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
			for _, winrate := range playerOnHeroWinrates {
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

		////////////////////////////////

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

		///////////////////////////////

		for _, winrate := range heroWinrates {
			rowData = append(rowData, fmt.Sprintf("%.2f", winrate.Winrate))
		}

		///////////////////////////////

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

		////////////////////////////////////////

		outputRow := strings.Join(rowData, ",")

		path, err := utils.GetPath("internal/ml/scripts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(path+"/row.txt", []byte(outputRow), 0644); err != nil {
			http.Error(w, "Failed to write to row.txt"+" err:"+err.Error(), http.StatusInternalServerError)
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

		file2, err := os.Open(path + "/prediction_result.json")
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

		//data["PredictionResult"] = resultData["prediction"]
		//data["PredictionProbability"] = resultData["probability"]
		prediction := resultData["prediction"].(bool)
		probability := resultData["probability"].(float64)

		var predictedTeam string
		var predictedProbability float64
		if probability < 0.5 {
			if prediction {
				predictedTeam = team2
				predictedProbability = (1 - probability) * 100
			} else {
				predictedTeam = team1
				predictedProbability = (1 - probability) * 100
			}
		} else {
			if prediction {
				predictedTeam = team1
				predictedProbability = probability * 100
			} else {
				predictedTeam = team2
				predictedProbability = (1 - probability) * 100
			}
		}

		data["PredictionResult"] = predictedTeam
		data["PredictionProbability"] = fmt.Sprintf("%.2f", predictedProbability)

	}

	if teams, err := hm.servicesManager.GetTeamsList(); err == nil {
		data["Team"] = teams
	}
	if heroList, err := hm.servicesManager.GetHeroesNameList(); err == nil {
		data["Hero"] = heroList
	}

	net.RenderTemplate(w, r, "prediction.html", data)
}
