package net

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func RedirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		if _, ok := session.Values["email"]; ok {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		email := session.Values["email"]

		if email == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SaveSession(w http.ResponseWriter, r *http.Request, email string) error {
	session, _ := store.Get(r, "session-name")
	session.Values["email"] = email
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return err
	}

	return nil
}
