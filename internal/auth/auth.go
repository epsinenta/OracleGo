package auth

import (
	"OracleGo/internal/db"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func findUser(email string) (db.User, error) {
	userManager, err := db.NewUsersDatabaseManager()
	if err != nil {
		log.Fatalf("Не удалось создать DataBaseManager: %v", err)
	}

	return userManager.GetUser(db.Email{email})
}

func ValidateUser(email string, password string) bool {
	gettenUser, parsErr := findUser(email)
	if parsErr != nil {

		if parsErr.Error() == "Пользователь не найден\n" {
			return false
		}
		return false
	}

	if (gettenUser.Email != db.Email{Value: email}) || !checkPassword(gettenUser.Password.GetValue(), password) {

		return false
	}
	return true
}
