package handlers

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/log"
	"OracleGo/internal/services"
	"OracleGo/internal/utils"
	"html/template"
	"net/http"

	"github.com/pkg/errors"
)

type HandlersManager struct {
	servicesManager *services.ServicesManager
	templates       *template.Template
	logger          *log.Logger
}

func NewHandlerManager(servicesManager *services.ServicesManager, logger *log.Logger) (*HandlersManager, error) {
	path, err := utils.GetPath("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	templ, err := template.ParseGlob(path)
	if err != nil {
		return nil, err
	}

	return &HandlersManager{servicesManager: servicesManager, templates: templ, logger: logger}, nil
}

func getBoolFromContext(r *http.Request, key interface{}) (bool, error) {
	flag := r.Context().Value(key)
	if flag == nil {
		return false, errors.New("getting value from context")
	}

	val, ok := flag.(bool)
	if !ok {
		return false, errors.New("type asserting from context")
	}

	return val, nil
}

// костыльная херня
func getClaimsFromContext(r *http.Request, key interface{}) (*entities.Claims, error) {
	claims := r.Context().Value(key)
	if claims == nil {
		return nil, errors.New("getting value from context")
	}

	val, ok := claims.(*entities.Claims)
	if !ok {
		return nil, errors.New("type asserting from context")
	}

	return val, nil
}
