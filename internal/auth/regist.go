package auth

func createUser(userManager UserManagerInterface, email string, password string) error {
	hashedPassword, hashErr := hashPassword(password)
	if hashErr != nil {
		return hashErr
	}
	return userManager.AddUsers([]Email{Email{email}}, []Password{Password{hashedPassword}})
}

func AddUser(userManager UserManagerInterface, email string, password string, confirmPassword string) bool {
	if password != confirmPassword {
		return false
	}
	_, parsErr := findUser(userManager, email)
	if parsErr != nil {

		if parsErr.Error() == "Пользователь не найден\n" {

			return createUser(userManager, email, password) == nil
		}

		return false
	}

	return false

}
