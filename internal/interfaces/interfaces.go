package interfaces

import (
	"OracleGo/internal/entities"
	"database/sql"
	"time"
)

type RedisManager interface {
	CacheData(key string, value interface{}, expiration time.Duration) error
	GetCachedData(key string, dest interface{}) error
	Ping() error
	GratefulStop() error
}
type DatabaseManager interface {
	Close() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	GratefulStop() error
}

type RepositoryManager interface {
	AddRows(tableName string, args map[string][]string) error
	BuildSQLQuery(tableName string, params []string, args map[string][]string) string
	DeleteRows(tableName string, args map[string][]string) error
	GetRows(tableName string, params []string, args map[string][]string) ([][]string, error)
}

type ServicesManager interface {
	AddUsers(emails []entities.Email, passwords []entities.Password) error
	GetAllHeroesWinrates() ([]entities.Winrate, error)
	GetHeroesCounterPicks(firstHeroes []entities.Hero, secondHeroes []entities.Hero) ([]entities.CounterRate, error)
	GetHeroesNameList() ([]entities.Hero, error)
	GetHeroesNameListPerPatch(patch string) ([]entities.Hero, error)
	GetHeroesWinrates(heroes []entities.Hero) ([]entities.Winrate, error)
	GetPlayerCountOnHero(players []entities.Player, heroes []entities.Hero) ([]entities.GamesCount, error)
	GetPlayerOnHeroWinrate(players []entities.Player, heroes []entities.Hero) ([]entities.PlayerWinrate, error)
	GetPlayersWinrate(players []entities.Player) ([]entities.PlayerWinrateCurPatch, error)
	GetTeamsList() ([]entities.Team, error)
	GetTeamsRoastersList() ([]entities.TeamRoaster, error)
	GetUser(email entities.Email) (entities.User, error)
}
