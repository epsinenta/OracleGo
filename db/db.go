package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
