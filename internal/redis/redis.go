package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisManager - структура для работы с Redis
type RedisManager struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	redisClient   *redis.Client
	ctx           context.Context
}

// NewRedisManager - конструктор для инициализации RedisManager
func NewRedisManager() *RedisManager {
	manager := &RedisManager{
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		ctx:           context.Background(),
	}

	manager.redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", manager.RedisHost, manager.RedisPort),
		Password: manager.RedisPassword, // пароль по умолчанию ""
		DB:       0,                     // используемая БД по умолчанию
	})

	// Проверка подключения
	if err := manager.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	/*
		err := manager.redisClient.FlushDB(manager.ctx).Err()
		if err != nil {
			log.Fatalf("Ошибка при очистке базы данных: %v", err)
		}
	*/
	return manager
}

// Ping - проверка подключения к Redis
func (r *RedisManager) Ping() error {
	_, err := r.redisClient.Ping(r.ctx).Result()
	return err
}

// CacheData - функция для кэширования данных
func (r *RedisManager) CacheData(key string, value interface{}, expiration time.Duration) error {
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

// GetCachedData - функция для получения кэшированных данных
func (r *RedisManager) GetCachedData(key string, dest interface{}) error {
	data, err := r.redisClient.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("данные по ключу '%s' не найдены", key)
	} else if err != nil {
		return fmt.Errorf("ошибка при получении данных из Redis: %v", err)
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("ошибка при десериализации данных: %v", err)
	}

	return nil
}

// Вспомогательная функция для получения переменных окружения с дефолтным значением
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
