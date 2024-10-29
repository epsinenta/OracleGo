package handlers

import (
	"OracleGo/internal/auth"
	_ "fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")

		isValidUser := auth.ValidateUser(email, inputPassword)
		if isValidUser {
			err := saveSession(w, r, email)
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

	renderTemplate(w, r, "login.html", data)
}
