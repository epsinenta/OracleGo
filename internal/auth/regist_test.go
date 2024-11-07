package auth

import (
	"errors"
	"testing"
)

type MockUserDatabaseManager struct {
	users map[string]string
}

func NewMockUserDatabaseManager() (*MockUserDatabaseManager, error) {
	return &MockUserDatabaseManager{users: make(map[string]string)}, nil
}

func (m *MockUserDatabaseManager) GetUser(email Email) (User, error) {
	password, exists := m.users[email.Value]
	if !exists {
		return User{}, errors.New("Пользователь не найден\n")
	}
	return User{
		Email:    email,
		Password: Password{Value: password},
	}, nil
}

func (m *MockUserDatabaseManager) AddUsers(emails []Email, passwords []Password) error {
	if len(emails) != len(passwords) {
		return errors.New("Несоответствие количества email и паролей")
	}
	for i, email := range emails {
		m.users[email.Value] = passwords[i].GetValue()
	}
	return nil
}

func TestCreateUser(t *testing.T) {
	manager, _ := NewMockUserDatabaseManager()
	email := "admin@admin.com"
	password := "admin"

	err := createUser(manager, email, password)
	if err != nil {
		t.Errorf("createUser returned an error: %v", err)
	}

	user, err := manager.GetUser(Email{Value: email})
	if err != nil {
		t.Errorf("User was not created: %v", err)
	}

	if !checkPassword(user.Password.GetValue(), password) {
		t.Errorf("Stored password does not match the original")
	}
}

func TestAddUser(t *testing.T) {
	tests := []struct {
		name            string
		email           string
		password        string
		confirmPassword string
		setup           func(manager *MockUserDatabaseManager)
		expected        bool
	}{
		{
			name:            "successful user creation",
			email:           "admin@admin.com",
			password:        "admin",
			confirmPassword: "admin",
			setup:           func(manager *MockUserDatabaseManager) {},
			expected:        true,
		},
		{
			name:            "passwords do not match",
			email:           "admin@admin.com",
			password:        "admin",
			confirmPassword: "notadmin",
			setup:           func(manager *MockUserDatabaseManager) {},
			expected:        false,
		},
		{
			name:            "user already exists",
			email:           "admin@admin.com",
			password:        "admin",
			confirmPassword: "admin",
			setup: func(manager *MockUserDatabaseManager) {
				hashedPassword, _ := hashPassword("admin")
				manager.users["admin@admin.com"] = hashedPassword
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, _ := NewMockUserDatabaseManager()
			tt.setup(manager)

			result := AddUser(manager, tt.email, tt.password, tt.confirmPassword)
			if result != tt.expected {
				t.Errorf("AddUser(%v, %v, %v) = %v; want %v", tt.email, tt.password, tt.confirmPassword, result, tt.expected)
			}
		})
	}
}
