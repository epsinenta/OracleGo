package db

import (
	"errors"
	_ "fmt"
	"log"

	_ "github.com/lib/pq"
)

type UsersDatabaseManager struct {
	dbManager DatabaseManager
}

func NewUsersDatabaseManager() (*UsersDatabaseManager, error) {
	dbManager, err := NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &UsersDatabaseManager{dbManager: *dbManager}, nil
}

func (usersDbManager *UsersDatabaseManager) GetUser(email Email) (User, error) {
	userRows, err := usersDbManager.dbManager.GetRows("users", []string{"email", "password"}, map[string][]string{"email": []string{email.Value}})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	if len(userRows) == 0 {
		return User{}, errors.New("Пользователь не найден\n")
	}
	var result User
	result = User{Email: Email{userRows[0][0]}, Password: Password{userRows[0][1]}}
	return result, nil
}

func (usersDbManager *UsersDatabaseManager) AddUsers(emails []Email, passwords []Password) error {
	emailsValues := ValuesFromAny(emails)
	passwordsValues := ValuesFromAny(passwords)

	err := usersDbManager.dbManager.AddRows("users", map[string][]string{"email": emailsValues, "password": passwordsValues})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	return nil
}
