package handlers

import (
	"OracleGo/internal/auth"
	_ "fmt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")
		isCreateUser := auth.AddUser(email, inputPassword, confirmPassword)
		if isCreateUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {
			data["ErrorMessage"] = "Incorrect email or password"
		}
	}

	renderTemplate(w, r, "register.html", data)
}
