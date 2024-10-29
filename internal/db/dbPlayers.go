package db

import (
	_ "fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type PlayersDatabaseManager struct {
	dbManager DatabaseManager
}

func NewPlayersDatabaseManager() (*PlayersDatabaseManager, error) {
	dbManager, err := NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &PlayersDatabaseManager{dbManager: *dbManager}, nil
}

func (playersDbManager *PlayersDatabaseManager) GetPlayerOnHeroWinrate(players []Player, heroes []Hero) ([]PlayerWinrate, error) {
	playerNames := ValuesFromAny(players)
	heroesNames := ValuesFromAny(heroes)
	winrateRows, err := playersDbManager.dbManager.GetRows("players_heroes_statistics", []string{"winrate", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []PlayerWinrate
	for _, row := range winrateRows {
		if winrate, convErr := strconv.ParseFloat(row[0], 64); err == nil {
			if convErr != nil {
				log.Fatalf("Не конвертировать string to float %v", err)
			}
			result = append(result, PlayerWinrate{Winrate: winrate, Player: Player{row[1]}, Hero: Hero{row[2]}})
		}
	}
	return result, nil
}

func (playersDbManager *PlayersDatabaseManager) GetPlayerCountOnHero(players []Player, heroes []Hero) ([]GamesCount, error) {
	playerNames := ValuesFromAny(players)
	heroesNames := ValuesFromAny(heroes)
	winrateRows, err := playersDbManager.dbManager.GetRows("players_heroes_statistics", []string{"count_of_matches", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []GamesCount
	for _, row := range winrateRows {
		if count, convErr := strconv.Atoi(row[0]); err == nil {
			if convErr != nil {
				log.Fatalf("Не конвертировать string to int %v", err)
			}
			result = append(result, GamesCount{Count: count, Player: Player{row[1]}, Hero: Hero{row[2]}})
		}
	}
	return result, nil
}
