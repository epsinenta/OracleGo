package handlers

import (
	"OracleGo/internal/interfaces"
	"OracleGo/internal/log"
)

type HandlersManager struct {
	servicesManager interfaces.ServicesManager
	logger          *log.Logger
}

func NewHandlerManager(servicesManager interfaces.ServicesManager, logger *log.Logger) *HandlersManager {
	return &HandlersManager{servicesManager: servicesManager, logger: logger}
}
