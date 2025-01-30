package redis

import (
	"os"
	"testing"
	"time"
)

// Тест для инициализации RedisManager
func TestNewRedisManager(t *testing.T) {
	// Устанавливаем переменные окружения для тестирования
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")

	_, err := NewRedisManager()
	if err != nil {
		t.Fatalf("Ошибка подключения к Redis: %v", err)
	}
}

// Тест для функции CacheData и GetCachedData
func TestCacheDataAndGetCachedData(t *testing.T) {
	manager, err := NewRedisManager()
	if err != nil {
		t.Fatal(err)
	}

	// Тестовые данные
	key := "testKey"
	value := map[string]string{"name": "Dota", "type": "Game"}
	expiration := 10 * time.Second

	// Сохраняем данные в кеш
	err = manager.CacheData(key, value, expiration)
	if err != nil {
		t.Fatalf("Ошибка кэширования данных: %v", err)
	}

	// Получаем данные из кеша
	var cachedValue map[string]string
	err = manager.GetCachedData(key, &cachedValue)
	if err != nil {
		t.Fatalf("Ошибка получения данных из кеша: %v", err)
	}

	// Проверяем, что значения совпадают
	if cachedValue["name"] != value["name"] || cachedValue["type"] != value["type"] {
		t.Fatalf("Полученные данные не совпадают с исходными. Ожидалось: %v, Получено: %v", value, cachedValue)
	}
}

// Тест на получение данных, когда ключа нет
func TestGetCachedDataNotFound(t *testing.T) {
	manager, err := NewRedisManager()
	if err != nil {
		t.Fatal(err)
	}

	var cachedValue map[string]string
	err = manager.GetCachedData("nonexistentKey", &cachedValue)
	if err == nil {
		t.Fatalf("Ожидалась ошибка при получении несуществующего ключа")
	}
}

// Тест для очистки кеша (не оч понятно зачем)

// func TestFlushDB(t *testing.T) {
// 	manager := NewRedisManager()

// 	// Устанавливаем тестовый ключ
// 	key := "flushTestKey"
// 	value := "someData"
// 	err := manager.CacheData(key, value, 5*time.Minute)
// 	if err != nil {
// 		t.Fatalf("Ошибка при кэшировании данных перед очисткой: %v", err)
// 	}

// 	// Очищаем базу данных
// 	err = manager.redisClient.FlushDB(manager.ctx).Err()
// 	if err != nil {
// 		t.Fatalf("Ошибка при очистке базы данных: %v", err)
// 	}

// 	// Проверяем, что данные удалены
// 	var cachedValue string
// 	err = manager.GetCachedData(key, &cachedValue)
// 	if err == nil {
// 		t.Fatalf("Ожидалась ошибка, т.к. данные должны были быть удалены")
// 	}
// }
