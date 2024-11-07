package net

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
  "os"
  "path/filepath"
)

var (
	store     = sessions.NewCookieStore([]byte("my-secret-key"))
	templates *template.Template
)

func init() {
    var path string
    if os.Getenv("GO_ENV") == "test" {
        path = filepath.Join("..", "..", "web", "templates", "*.html")
    } else {
        path = filepath.Join("web", "templates", "*.html")
    }
    templates = template.Must(template.ParseGlob(path))
}

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
