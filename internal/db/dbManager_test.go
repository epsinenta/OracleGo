package db

import (
	"os"
	"testing"
)

// setupTestDBManager инициализирует тестовый `DatabaseManager`.
func setupTestDBManager(t *testing.T) *DatabaseManager {
	t.Helper()

	// Устанавливаем переменную окружения для тестовой базы данных
	os.Setenv("GO_ENV", "test")
	dbManager, err := NewDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось создать DatabaseManager: %v", err)
	}
	return dbManager
}

// TestBuildSQLQuery проверяет генерацию SQL-запросов с различными параметрами.
func TestBuildSQLQuery(t *testing.T) {
	dbManager := setupTestDBManager(t)

	// Пример запроса с условиями
	tableName := "users"
	params := []string{"email", "password"}
	args := map[string][]string{"email": {"test@example.com"}}
	expectedQuery := "SELECT email, password FROM users WHERE email IN ('test@example.com')"

	query := dbManager.BuildSQLQuery(tableName, params, args)
	if query != expectedQuery {
		t.Errorf("Ожидался запрос %v, но получен %v", expectedQuery, query)
	}
}

// TestAddRows проверяет добавление строк в таблицу `users`.
func TestAddRows(t *testing.T) {
	dbManager := setupTestDBManager(t)
	defer dbManager.db.Close()

	// Добавление тестового пользователя
	args := map[string][]string{
		"email":    {"test@example.com"},
		"password": {"password"},
	}

	err := dbManager.AddRows("users", args)
	if err != nil {
		t.Fatalf("Ошибка при добавлении строки в таблицу users: %v", err)
	}

	// Удаление тестового пользователя после добавления
	defer dbManager.DeleteRows("users", map[string][]string{"email": {"test@example.com"}})

	// Проверяем, что данные добавлены
	rows, err := dbManager.GetRows("users", []string{"email", "password"}, map[string][]string{"email": {"test@example.com"}})
	if err != nil {
		t.Fatalf("Ошибка при получении строки из таблицы users: %v", err)
	}

	if len(rows) != 1 || rows[0][0] != "test@example.com" || rows[0][1] != "password" {
		t.Errorf("Ожидались данные email=test@example.com, password=password, но получены %v", rows)
	}
}

// TestGetRows проверяет выборку данных из таблицы `users`.
func TestGetRows(t *testing.T) {
	dbManager := setupTestDBManager(t)
	defer dbManager.db.Close()

	// Добавляем тестового пользователя
	err := dbManager.AddRows("users", map[string][]string{"email": {"test@example.com"}, "password": {"password"}})
	if err != nil {
		t.Fatalf("Ошибка при добавлении тестового пользователя: %v", err)
	}
	defer dbManager.DeleteRows("users", map[string][]string{"email": {"test@example.com"}})

	// Получаем данные пользователя
	rows, err := dbManager.GetRows("users", []string{"email", "password"}, map[string][]string{"email": {"test@example.com"}})
	if err != nil {
		t.Fatalf("Ошибка при получении строки из таблицы users: %v", err)
	}

	// Проверяем, что данные совпадают с ожидаемыми
	if len(rows) != 1 || rows[0][0] != "test@example.com" || rows[0][1] != "password" {
		t.Errorf("Ожидались данные email=test@example.com, password=password, но получены %v", rows)
	}
}

// TestDeleteRows проверяет удаление строк из таблицы `users`.
func TestDeleteRows(t *testing.T) {
	dbManager := setupTestDBManager(t)
	defer dbManager.db.Close()

	// Добавляем тестового пользователя
	err := dbManager.AddRows("users", map[string][]string{"email": {"delete_test@example.com"}, "password": {"password"}})
	if err != nil {
		t.Fatalf("Ошибка при добавлении тестового пользователя: %v", err)
	}

	// Удаляем тестового пользователя
	err = dbManager.DeleteRows("users", map[string][]string{"email": {"delete_test@example.com"}})
	if err != nil {
		t.Fatalf("Ошибка при удалении строки из таблицы users: %v", err)
	}

	// Проверяем, что данных больше нет
	rows, err := dbManager.GetRows("users", []string{"email"}, map[string][]string{"email": {"delete_test@example.com"}})
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса после удаления: %v", err)
	}

	if len(rows) != 0 {
		t.Errorf("Ожидалась пустая выборка после удаления, но получены данные: %v", rows)
	}
}
