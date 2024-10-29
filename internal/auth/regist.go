package auth

import (
	"OracleGo/internal/db"
	"log"
)

func createUser(email string, password string) error {
	userManager, err := db.NewUsersDatabaseManager()
	if err != nil {
		log.Fatalf("Не удалось создать DataBaseManager: %v", err)
	}
	hashedPassword, hashErr := hashPassword(password)
	if hashErr != nil {
		return hashErr
	}
	return userManager.AddUsers([]db.Email{db.Email{email}}, []db.Password{db.Password{hashedPassword}})
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
