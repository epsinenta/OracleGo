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

func (hm *HandlersManager) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		inputPassword := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		if inputPassword != confirmPassword {
			data["ErrorMessage"] = "passwords don't match"
		} else {
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				hm.logger.Log(err)
				return
			}

			_, err = hm.servicesManager.GetUser(entities.Email{Value: email})
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					err = hm.servicesManager.AddUsers([]entities.Email{{Value: email}}, []entities.Password{{Value: string(hashedPass)}})
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						hm.logger.Log(err)
						return
					}
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					hm.logger.Log(err)
					return
				}
			}

			data["ErrorMessage"] = "you are already registered"
		}
	}

	net.RenderTemplate(w, r, "register.html", data)
}
