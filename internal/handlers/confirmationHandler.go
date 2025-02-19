package handlers

import (
	"OracleGo/internal/entities"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func (hm *HandlersManager) ConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		sessionId := r.URL.Query().Get("sessionId")
		if len(sessionId) == 0 {
			err := errors.New("Empty session id")
			http.Error(w, err.Error(), http.StatusNotFound)
			hm.logger.Log(err)
			return
		}

		err := hm.servicesManager.Users.CompleteRegistration(sessionId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			hm.logger.Log(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     entities.RegistrationCookieName,
			Value:    sessionId,
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			Secure:   false, //поменять
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}
