package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Тест успешной авторизации
func TestLoginHandler_Success(t *testing.T) {
	hm := SetupTests(t, true)
	// Создание HTTP-запроса
	reqBody := strings.NewReader("email=test@example.com&password=password")
	req := httptest.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Вызов обработчика
	hm.LoginHandler(w, req)

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
	hm := SetupTests(t, false)

	reqBody := strings.NewReader("email=wrong@example.com&password=wrongpassword")
	req := httptest.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	hm.LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что обработчик не выполнил редирект и вернул страницу с ошибкой
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}

// Тест с методом GET, чтобы убедиться, что возвращается страница логина
func TestLoginHandler_GetMethod(t *testing.T) {
	hm := SetupTests(t, false)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	hm.LoginHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Проверка успешного рендеринга страницы логина
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}
