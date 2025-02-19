package statistics

import (
	"OracleGo/internal/entities"
	"fmt"

	_ "github.com/lib/pq"
)

func (sm *StatisticsManager) GetTeamsList() ([]entities.Team, error) {
	teamsRows, err := sm.repo.GetRows("teams_roasters", []string{"team_name"}, map[string][]string{})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.Team
	for _, row := range teamsRows {
		result = append(result, entities.Team{Value: row[0]})
	}
	return result, nil
}

func (sm *StatisticsManager) GetTeamsRoastersList() ([]entities.TeamRoaster, error) {
	teamsRows, err := sm.repo.GetRows("teams_roasters", []string{"*"}, map[string][]string{})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.TeamRoaster
	for _, row := range teamsRows {
		result = append(result, entities.TeamRoaster{
			Team: entities.Team{Value: row[0]},
			Players: []entities.Player{
				{Value: row[1]}, {Value: row[2]}, {Value: row[3]}, {Value: row[4]}, {Value: row[5]},
			},
		})
	}

	return result, nil
}
