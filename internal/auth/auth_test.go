package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "admin"
	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Errorf("hashPassword returned an error: %v", err)
	}

	if !checkPassword(hashedPassword, password) {
		t.Errorf("checkPassword returned false for correct password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "admin"
	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Errorf("hashPassword returned an error: %v", err)
	}

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"correct password", "admin", true},
		{"incorrect password", "wrongpassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPassword(hashedPassword, tt.password)
			if result != tt.expected {
				t.Errorf("checkPassword(%v) = %v; want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		expected bool
	}{
		{"valid user", "admin@admin.com", "admin", true},
		{"invalid password", "admin@admin.com", "wrongpassword", false},
		{"nonexistent user", "unknown@admin.com", "admin", false},
	}
	manager, _ := NewMockUserDatabaseManager()
	createUser(manager, "admin@admin.com", "admin")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUser(manager, tt.email, tt.password)
			if result != tt.expected {
				t.Errorf("ValidateUser(%v, %v) = %v; want %v", tt.email, tt.password, result, tt.expected)
			}
		})
	}
}
