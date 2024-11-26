package auth

import (
	"OracleGo/internal/db"
	"fmt"
	"testing"
)

func setupTestDatabaseManager() (*UsersDatabaseManager, error) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		return nil, fmt.Errorf("failed to setup test database manager: %v", err)
	}

	_ = dbManager.DeleteRows("users", map[string][]string{})

	return &UsersDatabaseManager{dbManager: *dbManager}, nil
}

func TestAddUsers(t *testing.T) {
	manager, err := setupTestDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}

	email := Email{Value: "testuser@example.com"}
	password := Password{Value: "hashedpassword123"}

	err = manager.AddUsers([]Email{email}, []Password{password})
	if err != nil {
		t.Errorf("AddUsers returned an error: %v", err)
	}

	user, err := manager.GetUser(email)
	if err != nil {

		t.Errorf("GetUser returned an error for existing user: %v", err)
	}
	if user.Email.GetValue() != email.GetValue() || user.Password.GetValue() != password.GetValue() {
		t.Errorf("User data does not match added values, got: %+v, want: %+v", user, User{Email: email, Password: password})
	}
}

func TestGetUser_NotFound(t *testing.T) {
	manager, err := setupTestDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}

	_, err = manager.GetUser(Email{Value: "nonexistent@example.com"})
	if err == nil || err.Error() != "Пользователь не найден\n" {
		t.Errorf("Expected 'Пользователь не найден' error, got: %v", err)
	}
}

func TestAddUsers_MultipleUsers(t *testing.T) {
	manager, err := setupTestDatabaseManager()
	if err != nil {
		t.Fatalf("Failed to set up test database manager: %v", err)
	}
	emails := []Email{{Value: "user1@example.com"}, {Value: "user2@example.com"}}
	passwords := []Password{{Value: "hashedpassword1"}, {Value: "hashedpassword2"}}
	err = manager.AddUsers(emails, passwords)
	if err != nil {
		t.Errorf("AddUsers returned an error: %v", err)
	}

	for i, email := range emails {
		user, err := manager.GetUser(email)
		if err != nil {
			t.Errorf("GetUser returned an error for added user %v: %v", email.GetValue(), err)
		}
		if user.Password.GetValue() != passwords[i].GetValue() {
			t.Errorf("User password does not match added value for %v, got: %s, want: %s", email.GetValue(), user.Password.GetValue(), passwords[i].GetValue())
		}
	}
}
