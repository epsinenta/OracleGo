package handlers

import (
	"OracleGo/internal/db"
	"OracleGo/internal/log"
	"OracleGo/internal/redis"
	"OracleGo/internal/repository"
	"OracleGo/internal/services"
	"testing"
)

func SetupTests(t *testing.T, setupUser bool) *HandlersManager {
	logger := log.New()

	db, err := db.NewDatabaseManager()
	if err != nil {
		t.Fatal(err)
	}

	redis, err := redis.NewRedisManager()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := db.GratefulStop(); err != nil {
			t.Error(err)
		}

		if err := redis.GratefulStop(); err != nil {
			t.Error(err)
		}
	})

	repo := repository.NewRepositoryManager(db, redis)

	svc, err := services.NewServiceManager(repo, redis)
	if err != nil {
		t.Fatal(err)
	}

	hm, err := NewHandlerManager(svc, logger)
	if err != nil {
		t.Fatal(err)
	}

	return hm
}
