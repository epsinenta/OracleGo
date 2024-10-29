package db

import (
	_ "fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type HeroesDatabaseManager struct {
	dbManager DatabaseManager
}

func NewHeroesDatabaseManager() (*HeroesDatabaseManager, error) {
	dbManager, err := NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &HeroesDatabaseManager{dbManager: *dbManager}, nil
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesNameList() ([]Hero, error) {
	return heroesDbManager.GetHeroesNameListPerPatch("7.35c")
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesNameListPerPatch(patch string) ([]Hero, error) {
	heroesRows, err := heroesDbManager.dbManager.GetRows("heroes_list", []string{"hero_name"}, map[string][]string{"patch": {patch}})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []Hero
	for _, row := range heroesRows {
		result = append(result, Hero{row[0]})
	}
	return result, nil
}

func (heroesDbManager *HeroesDatabaseManager) GetAllHeroesWinrates() ([]Winrate, error) {
	return heroesDbManager.GetHeroesWinrates(make([]Hero, 0))
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesWinrates(heroes []Hero) ([]Winrate, error) {
	heroesNames := ValuesFromAny(heroes)
	winratesRows, err := heroesDbManager.dbManager.GetRows("heroes_list", []string{"winrate", "hero_name"}, map[string][]string{"patch": {"7.35c"}, "hero_name": heroesNames})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var result []Winrate
	for _, row := range winratesRows {
		if winrate, convErr := strconv.ParseFloat(row[0], 64); err == nil {
			if convErr != nil {
				log.Fatalf("Не конвертировать string to float %v", err)
			}
			result = append(result, Winrate{Winrate: winrate, Hero: Hero{row[1]}})
		}

	}
	return result, nil
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesCounterPicks(firstHeroes []Hero, secondHeroes []Hero) ([]CounterRate, error) {
	firstHeroesNames := ValuesFromAny(firstHeroes)
	secondHeroesNames := ValuesFromAny(secondHeroes)
	var result []CounterRate
	for _, firstHero := range firstHeroesNames {
		multiFirstHero := strings.Split(strings.Repeat(firstHero+" ", len(secondHeroes)), " ")[:len(secondHeroes)]
		heroesRows, err := heroesDbManager.dbManager.GetRows("heroes_counters", []string{"counterrate", "first_hero_name", "second_hero_name"}, map[string][]string{"first_hero_name": multiFirstHero, "second_hero_name": secondHeroesNames})
		if err != nil {
			log.Fatalf("Не удалось провести запрос %v", err)
		}
		for _, row := range heroesRows {
			if counterrate, convErr := strconv.ParseFloat(row[0], 64); err == nil {
				if convErr != nil {
					log.Fatalf("Не конвертировать string to float %v", err)
				}
				result = append(result, CounterRate{CounterPick: counterrate, FirstHero: Hero{row[1]}, SecondHero: Hero{row[2]}})
			}

		}
	}

	return result, nil
}
