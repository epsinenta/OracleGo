package services

import (
	"OracleGo/internal/interfaces"
	"OracleGo/internal/services/statistics"
	"OracleGo/internal/services/users"
)

type ServicesManager struct {
	repo  interfaces.RepositoryManager
	cache interfaces.RedisManager

	Statistics *statistics.StatisticsManager
	Users      *users.UsersManager
}

func NewServiceManager(repo interfaces.RepositoryManager, cache interfaces.RedisManager) (*ServicesManager, error) {
	statistics := statistics.NewStatisticsManager(repo)
	users := users.NewUsersManager(repo, cache)
	return &ServicesManager{repo: repo, cache: cache, Statistics: statistics, Users: users}, nil
}
