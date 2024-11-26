package statistics

import (
	"OracleGo/internal/db"
	"fmt"
	_ "fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Hero struct {
	Value string
}

func (h Hero) GetValue() string {
	return h.Value
}

type Winrate struct {
	Hero    Hero
	Winrate float64
}

type CounterRate struct {
	FirstHero   Hero
	SecondHero  Hero
	CounterPick float64
}

type HeroesDatabaseManager struct {
	dbManager db.DatabaseManager
}

func NewHeroesDatabaseManager() (*HeroesDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &HeroesDatabaseManager{dbManager: *dbManager}, nil
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesNameList() ([]Hero, error) {
	return heroesDbManager.GetHeroesNameListPerPatch("7.35c")
}

func (heroesDbManager *HeroesDatabaseManager) getHeroesNameListPerPatchAsync(patch string, resultChan chan<- []Hero, errorChan chan<- error) {
	heroesRows, err := heroesDbManager.dbManager.GetRows("heroes_list", []string{"hero_name"}, map[string][]string{"patch": {patch}})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []Hero
	for _, row := range heroesRows {
		result = append(result, Hero{row[0]})
	}

	resultChan <- result
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesNameListPerPatch(patch string) ([]Hero, error) {
	resultChan := make(chan []Hero)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go heroesDbManager.getHeroesNameListPerPatchAsync(patch, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (heroesDbManager *HeroesDatabaseManager) GetAllHeroesWinrates() ([]Winrate, error) {
	return heroesDbManager.GetHeroesWinrates(make([]Hero, 0))
}

func (heroesDbManager *HeroesDatabaseManager) getHeroesWinratesAsync(heroes []Hero, resultChan chan<- []Winrate, errorChan chan<- error) {
	heroesNames := db.ValuesFromAny(heroes)
	params := map[string][]string{"patch": {"7.35c"}}
	if len(heroesNames) != 0 {
		params["hero_name"] = heroesNames
	}

	winratesRows, err := heroesDbManager.dbManager.GetRows("heroes_list", []string{"winrate", "hero_name"}, params)
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	var result []Winrate
	for _, row := range winratesRows {
		winrate, convErr := strconv.ParseFloat(row[0], 64)
		if convErr != nil {
			errorChan <- fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
			return
		}
		result = append(result, Winrate{Winrate: winrate, Hero: Hero{row[1]}})
	}

	resultChan <- result
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesWinrates(heroes []Hero) ([]Winrate, error) {
	resultChan := make(chan []Winrate)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go heroesDbManager.getHeroesWinratesAsync(heroes, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (heroesDbManager *HeroesDatabaseManager) getHeroesCounterPicksAsync(firstHeroes []Hero, secondHeroes []Hero, resultChan chan<- []CounterRate, errorChan chan<- error) {
	firstHeroesNames := db.ValuesFromAny(firstHeroes)
	secondHeroesNames := db.ValuesFromAny(secondHeroes)
	var result []CounterRate

	for _, firstHero := range firstHeroesNames {
		multiFirstHero := strings.Split(strings.Repeat(firstHero+" ", len(secondHeroes)), " ")[:len(secondHeroes)]
		heroesRows, err := heroesDbManager.dbManager.GetRows("heroes_counters", []string{"counterrate", "first_hero_name", "second_hero_name"}, map[string][]string{"first_hero_name": multiFirstHero, "second_hero_name": secondHeroesNames})
		if err != nil {
			errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
			return
		}

		for _, row := range heroesRows {
			counterrate, convErr := strconv.ParseFloat(row[0], 64)
			if convErr != nil {
				errorChan <- fmt.Errorf("не удалось конвертировать string to float: %w", convErr)
				return
			}
			result = append(result, CounterRate{CounterPick: counterrate, FirstHero: Hero{row[1]}, SecondHero: Hero{row[2]}})
		}
	}

	resultChan <- result
}

func (heroesDbManager *HeroesDatabaseManager) GetHeroesCounterPicks(firstHeroes []Hero, secondHeroes []Hero) ([]CounterRate, error) {
	resultChan := make(chan []CounterRate)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go heroesDbManager.getHeroesCounterPicksAsync(firstHeroes, secondHeroes, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
