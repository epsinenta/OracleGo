package auth

import (
	"OracleGo/internal/net"
	_ "fmt"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		userManager, err := NewUsersDatabaseManager()
		if err != nil {
			log.Fatalf("Не удалось создать DataBaseManager: %v", err)
		}

		isCreateUser := AddUser(userManager, email, inputPassword, confirmPassword)
		if isCreateUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {
			data["ErrorMessage"] = "Incorrect email or password"
		}
	}

	net.RenderTemplate(w, r, "register.html", data)
}
