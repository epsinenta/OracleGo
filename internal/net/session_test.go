package net

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectIfAuthenticated_RedirectsToProfile(t *testing.T) {
	// Создаем запрос и ответ
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Устанавливаем сессию с email для имитации аутентифицированного пользователя
	session, _ := store.Get(req, "session-name")
	session.Values["email"] = "test@example.com"
	session.Save(req, w)

	// Создаем тестовый обработчик, который вызывается, если не происходит редирект
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RedirectIfAuthenticated(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Проверяем, что произошел редирект на страницу профиля
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusSeeOther, resp.StatusCode)
	}
	if location := resp.Header.Get("Location"); location != "/profile" {
		t.Errorf("Ожидался редирект на /profile, но получен %v", location)
	}
}

func TestRedirectIfAuthenticated_AllowsAccessWhenUnauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Создаем тестовый обработчик, который должен быть вызван
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RedirectIfAuthenticated(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Проверяем, что запрос прошел без редиректа
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}

func TestSessionMiddleware_RedirectsToLoginWhenUnauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Создаем тестовый обработчик
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := SessionMiddleware(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Проверяем, что произошел редирект на страницу логина
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusSeeOther, resp.StatusCode)
	}
	if location := resp.Header.Get("Location"); location != "/login" {
		t.Errorf("Ожидался редирект на /login, но получен %v", location)
	}
}

func TestSessionMiddleware_AllowsAccessWhenAuthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Устанавливаем сессию с email для имитации аутентифицированного пользователя
	session, _ := store.Get(req, "session-name")
	session.Values["email"] = "test@example.com"
	session.Save(req, w)

	// Создаем тестовый обработчик, который должен быть вызван
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := SessionMiddleware(nextHandler)

	// Выполняем запрос
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Проверяем, что доступ разрешен, поскольку пользователь аутентифицирован
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}

func TestSaveSession_SavesSessionSuccessfully(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Сохраняем сессию с email
	err := SaveSession(w, req, "test@example.com")
	if err != nil {
		t.Fatalf("Ошибка сохранения сессии: %v", err)
	}

	// Проверяем, что сессия содержит email
	session, _ := store.Get(req, "session-name")
	email, ok := session.Values["email"].(string)
	if !ok || email != "test@example.com" {
		t.Errorf("Ожидался email 'test@example.com' в сессии, но получен %v", email)
	}
}
