package statistics

import (
	"OracleGo/internal/db"
	"fmt"
	"testing"
)

// Функция для настройки тестового окружения с очисткой таблицы перед каждым тестом
func setupTestHeroesDatabaseManager() (*HeroesDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, fmt.Errorf("failed to setup test database manager: %v", err)
	}

	return &HeroesDatabaseManager{dbManager: *dbManager}, nil
}

func TestGetHeroesNameListPerPatch(t *testing.T) {
	manager, err := setupTestHeroesDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}

	// Выполнение теста (без вставки данных)
	heroes, err := manager.GetHeroesNameListPerPatch("7.35c")
	if err != nil {
		t.Errorf("GetHeroesNameListPerPatch returned an error: %v", err)
		return
	}

	expectedHeroes := []Hero{{Value: "Anti-Mage"}, {Value: "Axe"}}
	if len(heroes) != len(expectedHeroes) {
		t.Errorf("Expected %d heroes, got %d", len(expectedHeroes), len(heroes))
	}
	for i, hero := range heroes {
		if hero.Value != expectedHeroes[i].Value {
			t.Errorf("Expected hero %s, got %s", expectedHeroes[i].Value, hero.Value)
		}
	}
}

func TestGetAllHeroesWinrates(t *testing.T) {
	manager, err := setupTestHeroesDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}

	// Выполнение теста (без вставки данных)
	winrates, err := manager.GetAllHeroesWinrates()
	if err != nil {
		t.Errorf("GetAllHeroesWinrates returned an error: %v", err)
		return
	}

	expectedWinrates := []Winrate{
		{Hero: Hero{Value: "Anti-Mage"}, Winrate: 53.2},
		{Hero: Hero{Value: "Axe"}, Winrate: 48.9},
	}
	if len(winrates) != len(expectedWinrates) {
		t.Errorf("Expected %d winrates, got %d", len(expectedWinrates), len(winrates))
	}
	for i, winrate := range winrates {
		if winrate.Hero.Value != expectedWinrates[i].Hero.Value || winrate.Winrate != expectedWinrates[i].Winrate {
			t.Errorf("Expected winrate %+v, got %+v", expectedWinrates[i], winrate)
		}
	}
}

func TestGetHeroesCounterPicks(t *testing.T) {
	manager, err := setupTestHeroesDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}

	// Выполнение теста (без вставки данных)
	firstHeroes := []Hero{{Value: "Anti-Mage"}}
	secondHeroes := []Hero{{Value: "Axe"}}
	counterPicks, err := manager.GetHeroesCounterPicks(firstHeroes, secondHeroes)
	if err != nil {
		t.Errorf("GetHeroesCounterPicks returned an error: %v", err)
		return
	}

	expectedCounterPicks := []CounterRate{
		{FirstHero: Hero{Value: "Anti-Mage"}, SecondHero: Hero{Value: "Axe"}, CounterPick: 55.1},
	}
	if len(counterPicks) != len(expectedCounterPicks) {

		t.Errorf("Expected %d counter picks, got %d", len(expectedCounterPicks), len(counterPicks))
	}
	for i, counterPick := range counterPicks {
		if counterPick.FirstHero.Value != expectedCounterPicks[i].FirstHero.Value || counterPick.SecondHero.Value != expectedCounterPicks[i].SecondHero.Value || counterPick.CounterPick != expectedCounterPicks[i].CounterPick {
			t.Errorf("Expected counter pick %+v, got %+v", expectedCounterPicks[i], counterPick)
		}
	}
}
