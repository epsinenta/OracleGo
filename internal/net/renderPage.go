package net

import (
	"html/template"
	"log"
	"net/http"

	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
)

var (
	store     = sessions.NewCookieStore([]byte("my-secret-key"))
	templates *template.Template
)

func GetTemplatePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		templatePath := filepath.Join(cwd, "web", "templates")

		if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
			return filepath.Join(templatePath, "*.html")
		}

		parentDir := filepath.Dir(cwd)
		if parentDir == cwd {
			log.Fatalf("Не удалось найти папку с шаблонами, начальный путь был: %s", filepath.Join(cwd, "web", "templates", "*.html"))
		}

		cwd = parentDir
	}
}

func init() {

	templatePath := GetTemplatePath()
	templates = template.Must(template.ParseGlob(templatePath))
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
