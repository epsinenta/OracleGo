package auth

import (
	"log"
)

func createUser(email string, password string) error {
	userManager, err := NewUsersDatabaseManager()
	if err != nil {
		log.Fatalf("Не удалось создать DataBaseManager: %v", err)
	}
	hashedPassword, hashErr := hashPassword(password)
	if hashErr != nil {
		return hashErr
	}
	return userManager.AddUsers([]Email{Email{email}}, []Password{Password{hashedPassword}})
}

func AddUser(email string, password string, confirmPassword string) bool {
	if password != confirmPassword {
		return false
	}
	_, parsErr := findUser(email)
	if parsErr != nil {

		if parsErr.Error() == "Пользователь не найден\n" {

			return createUser(email, password) == nil
		}

		return false
	}

	return false

}
