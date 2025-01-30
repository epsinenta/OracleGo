package handlers

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/net"
	"database/sql"
	_ "fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (hm *HandlersManager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")

		user, err := hm.servicesManager.GetUser(entities.Email{Value: email})
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				data["ErrorMessage"] = "Incorrect email or password"
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				hm.logger.Log(err)
				return
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password.Value), []byte(inputPassword))
		if err != nil {
			data["ErrorMessage"] = "Incorrect email or password"
			return
		} else {
			err = net.SaveSession(w, r, email)
			if err == nil {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				hm.logger.Log(err)
				return
			}
		}
	}

	net.RenderTemplate(w, r, "login.html", data)
}
