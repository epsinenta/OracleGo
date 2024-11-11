package statistics

import (
	"OracleGo/internal/db"
	"fmt"
	_ "fmt"

	_ "github.com/lib/pq"
)

type Team struct {
	Value string
}

func (t Team) GetValue() string {
	return t.Value
}

type TeamRoaster struct {
	Players []Player
	Team    Team
}

type TeamsDatabaseManager struct {
	dbManager db.DatabaseManager
}

func NewTeamsDatabaseManager() (*TeamsDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &TeamsDatabaseManager{dbManager: *dbManager}, nil
}
func (teamsDbManager *TeamsDatabaseManager) getTeamsListAsync(resultChan chan<- []Team, errorChan chan<- error) {
	teamsRows, err := teamsDbManager.dbManager.GetRows("teams_roasters", []string{"team_name"}, map[string][]string{})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []Team
	for _, row := range teamsRows {
		result = append(result, Team{row[0]})
	}

	// Отправка результата в канал
	resultChan <- result
}

func (teamsDbManager *TeamsDatabaseManager) GetTeamsList() ([]Team, error) {
	resultChan := make(chan []Team)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go teamsDbManager.getTeamsListAsync(resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (teamsDbManager *TeamsDatabaseManager) getTeamsRoastersListAsync(resultChan chan<- []TeamRoaster, errorChan chan<- error) {
	teamsRows, err := teamsDbManager.dbManager.GetRows("teams_roasters", []string{"*"}, map[string][]string{})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []TeamRoaster
	for _, row := range teamsRows {
		result = append(result, TeamRoaster{
			Team: Team{Value: row[0]},
			Players: []Player{
				{row[1]}, {row[2]}, {row[3]}, {row[4]}, {row[5]},
			},
		})
	}

	// Отправка результата в канал
	resultChan <- result
}

func (teamsDbManager *TeamsDatabaseManager) GetTeamsRoastersList() ([]TeamRoaster, error) {
	resultChan := make(chan []TeamRoaster)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go teamsDbManager.getTeamsRoastersListAsync(resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
