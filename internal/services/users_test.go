package services

import (
	"OracleGo/internal/entities"
	"testing"
)

func TestAddUsers(t *testing.T) {
	svc := SetupTests(t)

	email := entities.Email{Value: "testuser@example.com"}
	password := entities.Password{Value: "hashedpassword123"}

	err := svc.AddUsers([]entities.Email{email}, []entities.Password{password})
	if err != nil {
		t.Errorf("AddUsers returned an error: %v", err)
	}

	user, err := svc.GetUser(email)
	if err != nil {

		t.Errorf("GetUser returned an error for existing user: %v", err)
	}
	if user.Email.GetValue() != email.GetValue() || user.Password.GetValue() != password.GetValue() {
		t.Errorf("User data does not match added values, got: %+v, want: %+v", user, entities.User{Email: email, Password: password})
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc := SetupTests(t)

	_, err := svc.GetUser(entities.Email{Value: "nonexistent@example.com"})
	if err == nil || err.Error() != "Пользователь не найден\n" {
		t.Errorf("Expected 'Пользователь не найден' error, got: %v", err)
	}
}

func TestAddUsers_MultipleUsers(t *testing.T) {
	svc := SetupTests(t)

	emails := []entities.Email{{Value: "user1@example.com"}, {Value: "user2@example.com"}}
	passwords := []entities.Password{{Value: "hashedpassword1"}, {Value: "hashedpassword2"}}
	err := svc.AddUsers(emails, passwords)
	if err != nil {
		t.Errorf("AddUsers returned an error: %v", err)
	}

	for i, email := range emails {
		user, err := svc.GetUser(email)
		if err != nil {
			t.Errorf("GetUser returned an error for added user %v: %v", email.GetValue(), err)
		}
		if user.Password.GetValue() != passwords[i].GetValue() {
			t.Errorf("User password does not match added value for %v, got: %s, want: %s", email.GetValue(), user.Password.GetValue(), passwords[i].GetValue())
		}
	}
}
