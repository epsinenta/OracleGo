package users

import "OracleGo/internal/interfaces"

type UsersManager struct {
	repo  interfaces.RepositoryManager
	cache interfaces.RedisManager
}

func NewUsersManager(repo interfaces.RepositoryManager, cache interfaces.RedisManager) *UsersManager {
	return &UsersManager{repo: repo, cache: cache}
}
