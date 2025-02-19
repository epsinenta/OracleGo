package statistics

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/repository"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (sm *StatisticsManager) GetHeroesNameList() ([]entities.Hero, error) {
	return sm.GetHeroesNameListPerPatch("7.35c")
}

func (sm *StatisticsManager) GetHeroesNameListPerPatch(patch string) ([]entities.Hero, error) {
	heroesRows, err := sm.repo.GetRows("heroes_list", []string{"hero_name"}, map[string][]string{"patch": {patch}})
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.Hero
	for _, row := range heroesRows {
		result = append(result, entities.Hero{Value: row[0]})
	}

	return result, nil
}

func (sm *StatisticsManager) GetAllHeroesWinrates() ([]entities.Winrate, error) {
	return sm.GetHeroesWinrates(make([]entities.Hero, 0))
}

func (sm *StatisticsManager) GetHeroesWinrates(heroes []entities.Hero) ([]entities.Winrate, error) {
	heroesNames := repository.ValuesFromAny(heroes)
	params := map[string][]string{"patch": {"7.35c"}}
	if len(heroesNames) != 0 {
		params["hero_name"] = heroesNames
	}

	winratesRows, err := sm.repo.GetRows("heroes_list", []string{"winrate", "hero_name"}, params)
	if err != nil {
		return nil, fmt.Errorf("не удалось провести запрос: %w", err)
	}

	var result []entities.Winrate
	for _, row := range winratesRows {
		winrate, convErr := strconv.ParseFloat(row[0], 64)
		if convErr != nil {
			return nil, fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
		}
		result = append(result, entities.Winrate{Winrate: winrate, Hero: entities.Hero{Value: row[1]}})
	}

	return result, nil
}

func (sm *StatisticsManager) GetHeroesCounterPicks(firstHeroes []entities.Hero, secondHeroes []entities.Hero) ([]entities.CounterRate, error) {
	firstHeroesNames := repository.ValuesFromAny(firstHeroes)
	secondHeroesNames := repository.ValuesFromAny(secondHeroes)
	var result []entities.CounterRate

	for _, firstHero := range firstHeroesNames {
		multiFirstHero := strings.Split(strings.Repeat(firstHero+" ", len(secondHeroes)), " ")[:len(secondHeroes)]
		heroesRows, err := sm.repo.GetRows("heroes_counters", []string{"counterrate", "first_hero_name", "second_hero_name"}, map[string][]string{"first_hero_name": multiFirstHero, "second_hero_name": secondHeroesNames})
		if err != nil {
			return nil, fmt.Errorf("не удалось провести запрос: %w", err)
		}

		for _, row := range heroesRows {
			counterrate, convErr := strconv.ParseFloat(row[0], 64)
			if convErr != nil {
				return nil, fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
			}
			result = append(result, entities.CounterRate{CounterPick: counterrate, FirstHero: entities.Hero{Value: row[1]}, SecondHero: entities.Hero{Value: row[2]}})
		}
	}

	return result, nil
}
