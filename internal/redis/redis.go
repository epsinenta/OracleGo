package redis

import (
	"OracleGo/internal/interfaces"
	"OracleGo/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// redisManager - структура для работы с Redis
type redisManager struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	redisClient   *redis.Client
	ctx           context.Context
}

// NewRedisManager - конструктор для инициализации RedisManager
func NewRedisManager() (interfaces.RedisManager, error) {
	manager := &redisManager{
		RedisHost:     utils.GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     utils.GetEnv("REDIS_PORT", "6379"),
		RedisPassword: utils.GetEnv("REDIS_PASSWORD", ""),
		ctx:           context.Background(),
	}

	manager.redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", manager.RedisHost, manager.RedisPort),
		Password: manager.RedisPassword, // пароль по умолчанию ""
		DB:       0,                     // используемая БД по умолчанию
	})

	// Проверка подключения
	if err := manager.Ping(); err != nil {
		return nil, errors.Wrap(err, "Не удалось подключиться к Redis")
	}
	/*
		err := manager.redisClient.FlushDB(manager.ctx).Err()
		if err != nil {
			log.Fatalf("Ошибка при очистке базы данных: %v", err)
		}
	*/
	return manager, nil
}

// Ping - проверка подключения к Redis
func (r *redisManager) Ping() error {
	_, err := r.redisClient.Ping(r.ctx).Result()
	return err
}

// CacheData - функция для кэширования данных || json!!!
func (r *redisManager) CacheData(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации данных: %v", err)
	}

	err = r.redisClient.Set(r.ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка при сохранении данных в Redis: %v", err)
	}

	return nil
}

// GetCachedData - функция для получения кэшированных данных || json!!
func (r *redisManager) GetCachedData(key string, dest interface{}) error {
	data, err := r.redisClient.Get(r.ctx, key).Result()
	if err != nil {
		return errors.Wrap(err, "getting data from redis")
	}

	if dest != nil {
		err = json.Unmarshal([]byte(data), dest)
		if err != nil {
			return errors.Wrap(err, "unmarshaling json")
		}
	}

	return nil
}

func (r *redisManager) DeleteCachedData(keys []string) error {
	return errors.Wrap(r.redisClient.Del(r.ctx, keys...).Err(), "deleting from redis")
}

func (r *redisManager) GratefulStop() error {
	return r.redisClient.Close()
}
