package statistics

import (
	"testing"
)

func TestGetPlayerOnHeroWinrate(t *testing.T) {
	// Инициализация менеджера базы данных
	dbManager, err := NewPlayersDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось инициализировать менеджер базы данных: %v", err)
	}

	// Тестовые данные
	players := []Player{{Value: "Ame"}}
	heroes := []Hero{{Value: "Anti-Mage"}}

	// Выполнение функции
	result, err := dbManager.GetPlayerOnHeroWinrate(players, heroes)
	if err != nil {
		t.Fatalf("Ошибка при получении винрейта игрока на герое: %v", err)
	}

	// Проверка результатов
	if len(result) == 0 {
		t.Error("Ожидался непустой результат")
	} else {
		for _, winrate := range result {
			if winrate.Player.Value != "Ame" || winrate.Hero.Value != "Anti-Mage" || winrate.Winrate != 55.0 {
				t.Errorf("Неверный результат: ожидался {Player: 'Ame', Hero: 'Anti-Mage', Winrate: 55}, получено %+v", winrate)
			}
		}
	}
}
func TestGetPlayerCountOnHero(t *testing.T) {
	// Инициализация менеджера базы данных
	dbManager, err := NewPlayersDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось инициализировать менеджер базы данных: %v", err)
	}

	// Тестовые данные
	players := []Player{{Value: "Ame"}}
	heroes := []Hero{{Value: "Anti-Mage"}}

	// Выполнение функции
	result, err := dbManager.GetPlayerCountOnHero(players, heroes)
	if err != nil {
		t.Fatalf("Ошибка при получении количества матчей игрока на герое: %v", err)
	}

	// Проверка результатов
	if len(result) == 0 {
		t.Error("Ожидался непустой результат")
	} else {
		for _, count := range result {
			if count.Player.Value != "Ame" || count.Hero.Value != "Anti-Mage" || count.Count != 100 {
				t.Errorf("Неверный результат: ожидался {Player: 'Ame', Hero: 'Anti-Mage', Count: 100}, получено %+v", count)
			}
		}
	}
}
