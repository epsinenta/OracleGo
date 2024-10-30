package auth

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")

		isValidUser := ValidateUser(email, inputPassword)
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
