package auth

import (
	"OracleGo/internal/net"
	_ "fmt"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")
		userManager, err := NewUsersDatabaseManager()
		if err != nil {
			log.Fatalf("Не удалось создать DataBaseManager: %v", err)
		}
		isValidUser := ValidateUser(userManager, email, inputPassword)
		if isValidUser {
			err := net.SaveSession(w, r, email)
			if err == nil {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}
		} else {
			data["ErrorMessage"] = "Incorrect email or password"
		}
	}

	net.RenderTemplate(w, r, "login.html", data)
}
