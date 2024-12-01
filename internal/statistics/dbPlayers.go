package statistics

import (
	"OracleGo/internal/db"
	"fmt"
	_ "fmt"
	"strconv"

	_ "github.com/lib/pq"
)

type Player struct {
	Value string
}

func (p Player) GetValue() string {
	return p.Value
}

type PlayersDatabaseManager struct {
	dbManager db.DatabaseManager
}

type GamesCount struct {
	Player Player
	Hero   Hero
	Count  int
}

type PlayerWinrate struct {
	Player  Player
	Hero    Hero
	Winrate float64
}

func NewPlayersDatabaseManager() (*PlayersDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &PlayersDatabaseManager{dbManager: *dbManager}, nil
}

func (playersDbManager *PlayersDatabaseManager) getPlayerOnHeroWinrateAsync(players []Player, heroes []Hero, resultChan chan<- []PlayerWinrate, errorChan chan<- error) {
	playerNames := db.ValuesFromAny(players)
	heroesNames := db.ValuesFromAny(heroes)
	winrateRows, err := playersDbManager.dbManager.GetRows("players_heroes_statistic", []string{"winrate", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []PlayerWinrate
	for _, row := range winrateRows {
		winrate, convErr := strconv.ParseFloat(row[0], 64)
		if convErr != nil {
			errorChan <- fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
			return
		}
		result = append(result, PlayerWinrate{Winrate: winrate, Player: Player{row[1]}, Hero: Hero{row[2]}})
	}

	resultChan <- result
}

func (playersDbManager *PlayersDatabaseManager) GetPlayerOnHeroWinrate(players []Player, heroes []Hero) ([]PlayerWinrate, error) {
	resultChan := make(chan []PlayerWinrate)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go playersDbManager.getPlayerOnHeroWinrateAsync(players, heroes, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (playersDbManager *PlayersDatabaseManager) getPlayerCountOnHeroAsync(players []Player, heroes []Hero, resultChan chan<- []GamesCount, errorChan chan<- error) {
	playerNames := db.ValuesFromAny(players)
	heroesNames := db.ValuesFromAny(heroes)
	winrateRows, err := playersDbManager.dbManager.GetRows("players_heroes_statistic", []string{"count_of_matches", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []GamesCount
	for _, row := range winrateRows {
		count, convErr := strconv.Atoi(row[0])
		if convErr != nil {
			errorChan <- fmt.Errorf("не удалось конвертировать string to int: %w", convErr)
			return
		}
		result = append(result, GamesCount{Count: count, Player: Player{row[1]}, Hero: Hero{row[2]}})
	}

	resultChan <- result
}

func (playersDbManager *PlayersDatabaseManager) GetPlayerCountOnHero(players []Player, heroes []Hero) ([]GamesCount, error) {
	resultChan := make(chan []GamesCount)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go playersDbManager.getPlayerCountOnHeroAsync(players, heroes, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
