package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func findUser(userManager UserManagerInterface, email string) (User, error) {
	return userManager.GetUser(Email{email})
}

func ValidateUser(userManager UserManagerInterface, email string, password string) bool {
	gettenUser, parsErr := findUser(userManager, email)
	if parsErr != nil {

		if parsErr.Error() == "Пользователь не найден\n" {
			return false
		}
		return false
	}

	if (gettenUser.Email != Email{Value: email}) || !checkPassword(gettenUser.Password.GetValue(), password) {

		return false
	}
	return true
}
