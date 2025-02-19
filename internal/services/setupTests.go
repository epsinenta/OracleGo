package services

import (
	"OracleGo/internal/db"
	"OracleGo/internal/redis"
	"OracleGo/internal/repository"
	"testing"
)

func SetupTests(t *testing.T) *ServicesManager {
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

	svc, err := NewServiceManager(repo, redis)
	if err != nil {
		t.Fatal(err)
	}

	return svc
}
