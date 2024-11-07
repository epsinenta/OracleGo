package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestUser(t *testing.T) {
	manager, err := NewUsersDatabaseManager()
	if err != nil {
		t.Fatalf("Не удалось создать UsersDatabaseManager: %v", err)
	}

	AddUser(manager, "test@example.com", "password", "password")
}

// Тест успешной авторизации
func TestLoginHandler_Success(t *testing.T) {
	setupTestUser(t)
	// Создание HTTP-запроса
	reqBody := strings.NewReader("email=test@example.com&password=password")
	req := httptest.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Вызов обработчика
	LoginHandler(w, req)

	// Проверка результата
	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что был выполнен редирект на профиль
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusSeeOther, resp.StatusCode)
	}

	if loc := resp.Header.Get("Location"); loc != "/profile" {
		t.Errorf("Ожидался редирект на /profile, но получен %s", loc)
	}
}

// Тест с неверными данными авторизации
func TestLoginHandler_InvalidCredentials(t *testing.T) {
	reqBody := strings.NewReader("email=wrong@example.com&password=wrongpassword")
	req := httptest.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что обработчик не выполнил редирект и вернул страницу с ошибкой
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}

// Тест с методом GET, чтобы убедиться, что возвращается страница логина
func TestLoginHandler_GetMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Проверка успешного рендеринга страницы логина
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}
