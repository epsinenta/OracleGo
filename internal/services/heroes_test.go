package services

import (
	"OracleGo/internal/entities"
	"testing"
)

func TestGetHeroesNameListPerPatch(t *testing.T) {
	svc := SetupTests(t)

	// Выполнение теста (без вставки данных)
	heroes, err := svc.GetHeroesNameListPerPatch("7.35c")
	if err != nil {
		t.Errorf("GetHeroesNameListPerPatch returned an error: %v", err)
		return
	}

	expectedHeroes := []entities.Hero{{Value: "Anti-Mage"}, {Value: "Axe"}}
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
	svc := SetupTests(t)

	// Выполнение теста (без вставки данных)
	winrates, err := svc.GetAllHeroesWinrates()
	if err != nil {
		t.Errorf("GetAllHeroesWinrates returned an error: %v", err)
		return
	}

	expectedWinrates := []entities.Winrate{
		{Hero: entities.Hero{Value: "Anti-Mage"}, Winrate: 53.2},
		{Hero: entities.Hero{Value: "Axe"}, Winrate: 48.9},
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
	svc := SetupTests(t)

	// Выполнение теста (без вставки данных)
	firstHeroes := []entities.Hero{{Value: "Anti-Mage"}}
	secondHeroes := []entities.Hero{{Value: "Axe"}}
	counterPicks, err := svc.GetHeroesCounterPicks(firstHeroes, secondHeroes)
	if err != nil {
		t.Errorf("GetHeroesCounterPicks returned an error: %v", err)
		return
	}

	expectedCounterPicks := []entities.CounterRate{
		{FirstHero: entities.Hero{Value: "Anti-Mage"}, SecondHero: entities.Hero{Value: "Axe"}, CounterPick: 55.1},
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
