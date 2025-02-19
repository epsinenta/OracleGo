package statistics

import "OracleGo/internal/interfaces"

type StatisticsManager struct {
	repo interfaces.RepositoryManager
}

func NewStatisticsManager(repo interfaces.RepositoryManager) *StatisticsManager {
	return &StatisticsManager{repo: repo}
}
