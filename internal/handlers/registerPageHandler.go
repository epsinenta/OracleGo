package handlers

import (
	"OracleGo/internal/entities"
	_ "fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (hm *HandlersManager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
			email := r.FormValue("email")
			inputPassword := r.FormValue("password")
			confirmPassword := r.FormValue("confirm-password")

			if !hasMXRecord(email) {
				data[entities.ErrorMessageKey] = "wrong email"
			} else if inputPassword != confirmPassword {
				data[entities.ErrorMessageKey] = "passwords don't match"
			} else {
				sessionId, err := hm.servicesManager.Users.StartRegistration(entities.User{Email: email, Password: inputPassword})
				if err != nil {
					switch errors.Cause(err) {
					case entities.AlreadyRegisteredError:
						data[entities.ErrorMessageKey] = "you are already registered"
					case entities.SessionAlreadyStartedError:
						http.Redirect(w, r, "/profile-complete", http.StatusSeeOther)
						return
					default:
						http.Error(w, err.Error(), http.StatusInternalServerError)
						hm.logger.Log(err)
						return
					}
				} else {
					http.SetCookie(w, &http.Cookie{
						Name:     entities.RegistrationCookieName,
						Value:    sessionId,
						HttpOnly: true,
						Secure:   false, //поменять
						SameSite: http.SameSiteStrictMode,
						Path:     "/",
					})

					http.Redirect(w, r, "/profile-complete", http.StatusSeeOther)
					return
				}
			}
		}
	} else {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	if err := hm.templates.ExecuteTemplate(w, "register.html", data); err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func hasMXRecord(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	mxRecords, err := net.LookupMX(parts[1])
	return err == nil && len(mxRecords) > 0
}
