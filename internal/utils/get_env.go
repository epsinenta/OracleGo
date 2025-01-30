package utils

import "os"

// Вспомогательная функция для получения переменных окружения с дефолтным значением
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
