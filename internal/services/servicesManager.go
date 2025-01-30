package services

import "OracleGo/internal/interfaces"

type servicesManager struct {
	repo interfaces.RepositoryManager
}

func NewServiceManager(repo interfaces.RepositoryManager) interfaces.ServicesManager {
	return &servicesManager{repo: repo}
}
