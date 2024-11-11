package statistics

import (
	"testing"
)

func TestGetTeamsList(t *testing.T) {
	// Инициализация менеджера базы данных команд
	teamsDbManager, err := NewTeamsDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось создать менеджер базы данных: %v", err)
	}

	// Получение списка команд
	teams, err := teamsDbManager.GetTeamsList()
	if err != nil {
		t.Fatalf("Ошибка при выполнении GetTeamsList: %v", err)
	}

	// Проверка наличия тестовой команды
	expectedTeam := "Nouns"
	found := false
	for _, team := range teams {
		if team.Value == expectedTeam {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Команда %s не найдена в результате", expectedTeam)
	}
}

func TestGetTeamsRoastersList(t *testing.T) {
	// Инициализация менеджера базы данных команд
	teamsDbManager, err := NewTeamsDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось создать менеджер базы данных: %v", err)
	}

	// Получение списка составов команд
	roasters, err := teamsDbManager.GetTeamsRoastersList()
	if err != nil {
		t.Fatalf("Ошибка при выполнении GetTeamsRoastersList: %v", err)
	}

	// Проверка наличия тестового состава команды
	expectedTeam := "Nouns"
	expectedPlayers := []string{"A", "B", "C", "D", "E"}

	found := false
	for _, roaster := range roasters {
		if roaster.Team.Value == expectedTeam {
			found = true
			for i, player := range roaster.Players {
				if player.GetValue() != expectedPlayers[i] {
					t.Errorf("Ожидаемый игрок %s, но найден %s", expectedPlayers[i], player.GetValue())
				}
			}
			break
		}
	}

	if !found {
		t.Errorf("Состав команды %s не найден в результате", expectedTeam)
	}
}
