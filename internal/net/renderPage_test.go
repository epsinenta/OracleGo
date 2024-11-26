package net

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Функция для загрузки шаблонов в тестовом окружении
func loadTemplatesForTesting() {
	path := filepath.Join("..", "..", "web", "templates", "*.html")
	templates = template.Must(template.ParseGlob(path))
}

// Тест функции RenderTemplate, когда пользователь авторизован
func TestRenderTemplate_WithLoggedInUser(t *testing.T) {
	// Установка переменной окружения для тестовой базы шаблонов
	os.Setenv("GO_ENV", "test")
	defer os.Unsetenv("GO_ENV") // Очистка переменной окружения после теста

	// Загрузка шаблонов для тестовой среды
	loadTemplatesForTesting()

	// Создание тестового запроса и респондера
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Установка сессии с email для имитации авторизованного пользователя
	session, _ := store.Get(req, "session-name")
	session.Values["email"] = "test@example.com"
	session.Save(req, w)

	// Вызов RenderTemplate с авторизованным пользователем
	RenderTemplate(w, req, "example.html", map[string]interface{}{"Title": "Test Title"})

	resp := w.Result()
	defer resp.Body.Close()

	// Проверка, что статус ответа 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Проверка, что тело ответа содержит сообщение для авторизованного пользователя
	body := w.Body.String()
	if !strings.Contains(body, "Welcome, you are logged in!") {
		t.Errorf("Expected response to contain 'Welcome, you are logged in!', but got %s", body)
	}
}

// Тест функции RenderTemplate, когда пользователь не авторизован
func TestRenderTemplate_WithAnonymousUser(t *testing.T) {
	// Установка переменной окружения для тестовой базы шаблонов
	os.Setenv("GO_ENV", "test")
	defer os.Unsetenv("GO_ENV") // Очистка переменной окружения после теста

	// Загрузка шаблонов для тестовой среды
	loadTemplatesForTesting()

	// Создание тестового запроса и респондера
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Вызов RenderTemplate с анонимным пользователем (без сессии)
	RenderTemplate(w, req, "example.html", map[string]interface{}{"Title": "Test Title"})

	resp := w.Result()
	defer resp.Body.Close()

	// Проверка, что статус ответа 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Проверка, что тело ответа содержит сообщение для неавторизованного пользователя
	body := w.Body.String()
	if !strings.Contains(body, "Please log in to continue.") {
		t.Errorf("Expected response to contain 'Please log in to continue.', but got %s", body)
	}
}
