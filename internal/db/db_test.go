package db

import (
	"testing"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "1q2ws3edc4r"
	dbname   = "DotaTest" // тестовая база данных
)

// setupTestDB подключается к тестовой базе данных для использования в тестах.
func setupTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := NewDB(host, port, user, password, dbname)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой базе данных: %v", err)
	}
	return db
}

// TestNewDB_Success проверяет успешное подключение к базе данных.
func TestNewDB_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
}

// TestNewDB_InvalidCredentials проверяет, что при неверных данных подключения функция возвращает ошибку.
func TestNewDB_InvalidCredentials(t *testing.T) {
	_, err := NewDB(host, port, "wrong_user", "wrong_password", dbname)
	if err == nil {
		t.Error("Ожидалась ошибка при подключении с неверными учетными данными, но её не произошло")
	}
}

// TestDB_AddUser проверяет добавление пользователя в таблицу `users`.
func TestDB_AddUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Добавление пользователя в таблицу `users`.
	_, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "test@example.com", "password")
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса INSERT: %v", err)
	}

	// Удаление добавленного пользователя для чистоты тестов.
	_, err = db.Exec("DELETE FROM users WHERE email = $1", "test@example.com")
	if err != nil {
		t.Errorf("Ошибка удаления тестового пользователя: %v", err)
	}
}

// TestDB_GetUser проверяет выборку данных из таблицы `users`.
func TestDB_GetUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Добавляем тестового пользователя.
	_, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "test@example.com", "password")
	if err != nil {
		t.Fatalf("Ошибка добавления тестового пользователя: %v", err)
	}
	defer db.Exec("DELETE FROM users WHERE email = $1", "test@example.com") // Удаляем пользователя после теста.

	// Выполнение запроса SELECT.
	rows, err := db.Query("SELECT email, password FROM users WHERE email = $1", "test@example.com")
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса SELECT: %v", err)
	}
	defer rows.Close()

	// Проверяем, что результаты получены.
	if !rows.Next() {
		t.Error("Ожидался хотя бы один результат, но строки отсутствуют")
	}

	// Проверяем значения полей email и password.
	var email, password string
	if err := rows.Scan(&email, &password); err != nil {
		t.Errorf("Ошибка при чтении результата: %v", err)
	}
	if email != "test@example.com" || password != "password" {
		t.Errorf("Ожидались email 'test@example.com' и password 'password', но получены %s и %s", email, password)
	}
}

// TestDB_Close проверяет корректное закрытие соединения с базой данных.
func TestDB_Close(t *testing.T) {
	db := setupTestDB(t)

	// Закрытие подключения.
	err := db.Close()
	if err != nil {
		t.Errorf("Ошибка при закрытии подключения: %v", err)
	}

	// Проверка, что соединение закрыто.
	if err := db.conn.Ping(); err == nil {
		t.Error("Ожидалась ошибка при пинге закрытого подключения, но соединение активно")
	}
}
