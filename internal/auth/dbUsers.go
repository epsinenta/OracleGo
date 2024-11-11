package auth

import (
	"OracleGo/internal/db"
	"errors"
	"fmt"
	_ "fmt"

	_ "github.com/lib/pq"
)

type Email struct {
	Value string
}

func (e Email) GetValue() string {
	return e.Value
}

type Password struct {
	Value string
}

func (p Password) GetValue() string {
	return p.Value
}

type User struct {
	Email    Email
	Password Password
}

type UserManagerInterface interface {
	GetUser(email Email) (User, error)
	AddUsers(emails []Email, passwords []Password) error
}

type UsersDatabaseManager struct {
	dbManager db.DatabaseManager
}

func NewUsersDatabaseManager() (*UsersDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, err
	}
	return &UsersDatabaseManager{dbManager: *dbManager}, nil
}

func (usersDbManager *UsersDatabaseManager) getUserAsync(email Email, resultChan chan<- User, errorChan chan<- error) {
	userRows, err := usersDbManager.dbManager.GetRows("users", []string{"email", "password"}, map[string][]string{"email": {email.Value}})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}
	if len(userRows) == 0 {
		errorChan <- errors.New("Пользователь не найден\n")
		return
	}

	result := User{
		Email:    Email{userRows[0][0]},
		Password: Password{userRows[0][1]},
	}
	// Отправка результата в канал
	resultChan <- result
}

func (usersDbManager *UsersDatabaseManager) GetUser(email Email) (User, error) {
	resultChan := make(chan User)
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей каналов
	go usersDbManager.getUserAsync(email, resultChan, errorChan)

	// Ожидание результата или ошибки
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return User{}, err
	}
}

func (usersDbManager *UsersDatabaseManager) addUsersAsync(emails []Email, passwords []Password, errorChan chan<- error) {
	emailsValues := db.ValuesFromAny(emails)
	passwordsValues := db.ValuesFromAny(passwords)

	err := usersDbManager.dbManager.AddRows("users", map[string][]string{"email": emailsValues, "password": passwordsValues})
	if err != nil {
		errorChan <- fmt.Errorf("не удалось провести запрос: %w", err)
		return
	}

	// Отправка nil в errorChan для обозначения успешного завершения
	errorChan <- nil
}

func (usersDbManager *UsersDatabaseManager) AddUsers(emails []Email, passwords []Password) error {
	errorChan := make(chan error)

	// Запуск асинхронной функции с передачей канала ошибок
	go usersDbManager.addUsersAsync(emails, passwords, errorChan)

	// Ожидание результата или ошибки
	if err := <-errorChan; err != nil {
		return err
	}
	return nil
}
