package db

import (
	"OracleGo/internal/interfaces"
	"OracleGo/internal/utils"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type databaseManager struct {
	conn *sql.DB
}

func NewDatabaseManager() (interfaces.DatabaseManager, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		utils.GetEnv("DB_HOST", "postgres"),
		utils.GetEnv("DB_PORT", "5432"),
		utils.GetEnv("DB_USER", "postgres"),
		utils.GetEnv("DB_PASSWORD", "1q2ws3edc4r"),
		utils.GetEnv("DB_NAME", "Dota"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &databaseManager{conn: db}, nil
}

func (db *databaseManager) Close() error {
	return db.conn.Close()
}

func (db *databaseManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *databaseManager) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *databaseManager) GratefulStop() error {
	return db.conn.Close()
}
