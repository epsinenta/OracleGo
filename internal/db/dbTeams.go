package db

import (
	_ "fmt"
	"log"

	_ "github.com/lib/pq"
)

type TeamsDatabaseManager struct {
	dbManager DatabaseManager
}

func NewTeamsDatabaseManager() (*TeamsDatabaseManager, error) {
	dbManager, err := NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &TeamsDatabaseManager{dbManager: *dbManager}, nil
}

func (teamsDbManager *TeamsDatabaseManager) GetTeamsList() ([]Team, error) {
	teamsRows, err := teamsDbManager.dbManager.GetRows("teams_roasters", []string{"team_name"}, map[string][]string{})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []Team
	for _, row := range teamsRows {
		result = append(result, Team{row[0]})
	}
	return result, nil
}

func (teamsDbManager *TeamsDatabaseManager) GetTeamsRoastersList() ([]TeamRoaster, error) {
	teamsRows, err := teamsDbManager.dbManager.GetRows("teams_roasters", []string{"*"}, map[string][]string{})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []TeamRoaster
	for _, row := range teamsRows {
		result = append(result, TeamRoaster{Team: Team{Name: row[0]}, Players: []Player{{row[1]}, {row[2]}, {row[3]}, {row[4]}, {row[5]}}})
	}
	return result, nil
}
