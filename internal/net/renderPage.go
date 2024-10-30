package net

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	store     = sessions.NewCookieStore([]byte("my-secret-key"))
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data map[string]interface{}) {
	session, _ := store.Get(r, "session-name")
	_, isLoggedIn := session.Values["email"]

	if data == nil {
		data = make(map[string]interface{})
	}
	data["IsLoggedIn"] = isLoggedIn
	if err := templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
