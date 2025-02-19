package handlers

import (
	"OracleGo/internal/entities"
	"fmt"
	"net/http"
)

func (hm *HandlersManager) ProfileCompleteHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	isInReg, err := getBoolFromContext(r, entities.RegistrationStatusContextKey{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hm.logger.Log(err)
		return
	}

	if !isInReg {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	claims, err := getClaimsFromContext(r, entities.AuthClaimsContextKey{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		hm.logger.Log(err)
		return
	}

	if r.Method == http.MethodPost {
		action := r.URL.Query().Get("action")
		if action == "" {
			http.Error(w, "wrong query", http.StatusBadRequest)
			return
		}

		switch action {
		case "resend":
			err := hm.servicesManager.Users.ResendEmail(claims.SessionId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				hm.logger.Log(err)
				return
			}
		case "logout":
			fmt.Println("logout")
		default:
			http.Error(w, "wrong query", http.StatusBadRequest)
			return
		}
	}

	if err := hm.templates.ExecuteTemplate(w, "profile-complete.html", data); err != nil {
		hm.logger.Log(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
