package services

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/repository"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (sm *servicesManager) GetUser(email entities.Email) (entities.User, error) {
	userRows, err := sm.repo.GetRows("users", []string{"email", "password"}, map[string][]string{"email": {email.Value}})
	if err != nil {
		return entities.User{}, fmt.Errorf("не удалось провести запрос: %w", err)
	}
	if len(userRows) == 0 {
		return entities.User{}, sql.ErrNoRows
	}

	result := entities.User{
		Email:    entities.Email{Value: userRows[0][0]},
		Password: entities.Password{Value: userRows[0][1]},
	}

	return result, nil
}

func (sm *servicesManager) AddUsers(emails []entities.Email, passwords []entities.Password) error {
	emailsValues := repository.ValuesFromAny(emails)
	passwordsValues := repository.ValuesFromAny(passwords)

	err := sm.repo.AddRows("users", map[string][]string{"email": emailsValues, "password": passwordsValues})
	if err != nil {
		return fmt.Errorf("не удалось провести запрос: %w", err)
	}

	return nil
}
