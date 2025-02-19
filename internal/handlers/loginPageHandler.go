package handlers

import (
	"OracleGo/internal/entities"
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (hm *HandlersManager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	isLoggedIn, err := getBoolFromContext(r, entities.AuthStatusContextKey{})
	if err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isInRegistration, err := getBoolFromContext(r, entities.RegistrationStatusContextKey{})
	if err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isInRegistration {
		http.Redirect(w, r, "/profile-complete", http.StatusSeeOther)
		return
	}

	data[entities.IsLoggedInKey] = isLoggedIn
	if !isLoggedIn {
		if r.Method == http.MethodPost {
			inputEmail := r.FormValue("email")
			inputPassword := r.FormValue("password")

			token, err := hm.servicesManager.Users.Login(entities.User{Email: inputEmail, Password: inputPassword})
			if err != nil {
				switch errors.Cause(err) {
				case sql.ErrNoRows, bcrypt.ErrMismatchedHashAndPassword:
					data[entities.ErrorMessageKey] = "incorrect email or password"
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
					hm.logger.Log(err)
					return
				}
			}

			http.SetCookie(w, &http.Cookie{
				Name:     entities.RefreshTokenCookieName,
				Value:    token,
				HttpOnly: true,
				Secure:   false, //поменять
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
	} else {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	if err := hm.templates.ExecuteTemplate(w, "login.html", data); err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
