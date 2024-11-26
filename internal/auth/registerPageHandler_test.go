package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandler_Success(t *testing.T) {
	// Подготовка тестовых данных
	setupTestUser(t)

	// Создание HTTP-запроса для регистрации
	reqBody := strings.NewReader("email=newuser@example.com&password=password&confirm-password=password")
	req := httptest.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Вызов обработчика
	RegisterHandler(w, req)

	// Получаем результат
	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем статус ответа (должен быть редирект на /login)
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusSeeOther, resp.StatusCode)
	}

	// Проверяем редирект на /login
	loc := resp.Header.Get("Location")
	if loc != "/login" {
		t.Errorf("Ожидался редирект на /login, но получен %s", loc)
	}
}
func TestRegisterHandler_PasswordsDoNotMatch(t *testing.T) {
	// Создание HTTP-запроса для регистрации с разными паролями
	reqBody := strings.NewReader("email=newuser@example.com&password=password&confirm-password=wrongpassword")
	req := httptest.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Вызов обработчика
	RegisterHandler(w, req)

	// Получаем результат
	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что не произошло редиректа, а ошибка отображается на странице
	if resp.StatusCode == http.StatusSeeOther {
		t.Errorf("Ожидался статус не редиректа, но получен %v", resp.StatusCode)
	}

	// Проверяем, что сообщение об ошибке присутствует
	bodyBytes, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(bodyBytes), "Incorrect email or password") {
		t.Errorf("Ожидалось сообщение об ошибке, но его нет")
	}
}
func TestRegisterHandler_EmailAlreadyExists(t *testing.T) {
	// Добавляем пользователя с таким email для симуляции ситуации, когда email уже существует
	setupTestUser(t)

	// Создание HTTP-запроса для регистрации с уже существующим email
	reqBody := strings.NewReader("email=test@example.com&password=password&confirm-password=password")
	req := httptest.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Вызов обработчика
	RegisterHandler(w, req)

	// Получаем результат
	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что не произошло редиректа
	if resp.StatusCode == http.StatusSeeOther {
		t.Errorf("Ожидался статус не редиректа, но получен %v", resp.StatusCode)
	}

	// Проверяем, что ошибка связана с уже существующим email
	bodyBytes, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(bodyBytes), "Incorrect email or password") {
		t.Errorf("Ожидалось сообщение об ошибке, но его нет")
	}
}
func TestRegisterHandler_GetRequest(t *testing.T) {
	// Создание GET-запроса на страницу регистрации
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	w := httptest.NewRecorder()

	// Вызов обработчика
	RegisterHandler(w, req)

	// Получаем результат
	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем, что запрос вернул статус 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}

	// Дополнительно можно проверить, что тело ответа содержит форму регистрации
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Проверяем, что в теле ответа есть поля для ввода email и пароля
	if !strings.Contains(bodyStr, "name=\"email\"") {
		t.Errorf("Ожидалось поле для ввода email, но его нет в теле ответа")
	}

	if !strings.Contains(bodyStr, "name=\"password\"") {
		t.Errorf("Ожидалось поле для ввода пароля, но его нет в теле ответа")
	}

	if !strings.Contains(bodyStr, "name=\"confirm-password\"") {
		t.Errorf("Ожидалось поле для ввода подтверждения пароля, но его нет в теле ответа")
	}
}
