package services

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/repository"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

func (sm *servicesManager) GetPlayerOnHeroWinrate(players []entities.Player, heroes []entities.Hero) ([]entities.PlayerWinrate, error) {
	playerNames := repository.ValuesFromAny(players)
	heroesNames := repository.ValuesFromAny(heroes)
	winrateRows, err := sm.repo.GetRows("players_heroes_statistic", []string{"winrate", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.PlayerWinrate
	for _, row := range winrateRows {
		winrate, convErr := strconv.ParseFloat(row[0], 64)
		if convErr != nil {
			return nil, fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
		}
		result = append(result, entities.PlayerWinrate{Winrate: winrate, Player: entities.Player{Value: row[1]}, Hero: entities.Hero{Value: row[2]}})
	}

	return result, nil
}

func (sm *servicesManager) GetPlayersWinrate(players []entities.Player) ([]entities.PlayerWinrateCurPatch, error) {
	playerNames := repository.ValuesFromAny(players)
	winrateRows, err := sm.repo.GetRows("pro_players_list", []string{"winrate", "player_name"}, map[string][]string{"player_name": playerNames, "patch": {"7.35c"}})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.PlayerWinrateCurPatch
	for _, row := range winrateRows {
		winrate, convErr := strconv.ParseFloat(row[0], 64)
		if convErr != nil {
			return nil, fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
		}
		result = append(result, entities.PlayerWinrateCurPatch{Winrate: winrate, Player: entities.Player{Value: row[1]}})
	}

	return result, nil
}

func (sm *servicesManager) GetPlayerCountOnHero(players []entities.Player, heroes []entities.Hero) ([]entities.GamesCount, error) {
	playerNames := repository.ValuesFromAny(players)
	heroesNames := repository.ValuesFromAny(heroes)
	winrateRows, err := sm.repo.GetRows("players_heroes_statistic", []string{"count_of_matches", "player_name", "hero_name"}, map[string][]string{"player_name": playerNames, "hero_name": heroesNames})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.GamesCount
	for _, row := range winrateRows {
		count, convErr := strconv.Atoi(row[0])
		if convErr != nil {
			return nil, fmt.Errorf("не удалось конвертировать string to int: %w", convErr)
		}
		result = append(result, entities.GamesCount{Count: count, Player: entities.Player{Value: row[1]}, Hero: entities.Hero{Value: row[2]}})
	}

	return result, nil
}
