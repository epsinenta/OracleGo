package db

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type DatabaseManager struct {
	db DB
}

func NewDatabaseManager() (*DatabaseManager, error) {
	db, err := NewDB("localhost", "5432", "postgres", "1q2ws3edc4r", "Dota")
	if err != nil {
		return nil, err
	}
	return &DatabaseManager{db: *db}, nil
}

func (dbManager *DatabaseManager) BuildSQLQuery(tableName string, params []string, args map[string][]string) string {
	query := "SELECT " + strings.Join(params, ", ") + " FROM " + tableName

	if len(args) > 0 {
		var whereClauses []string
		for key, values := range args {
			if len(values) > 0 {
				placeholders := make([]string, len(values))
				for i := range values {
					placeholders[i] = fmt.Sprintf("'%s'", values[i])
				}
				whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
			}
		}
		if len(whereClauses) > 0 {
			query += " WHERE " + strings.Join(whereClauses, " AND ")
		}
	}

	return query
}

func (dbManager *DatabaseManager) GetRows(tableName string, params []string, args map[string][]string) ([][]string, error) {
	query := dbManager.BuildSQLQuery(tableName, params, args)
	rows, err := dbManager.db.Query(query)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var result [][]string

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				row[i] = fmt.Sprintf("%v", val)
			} else {
				row[i] = "NULL"
			}
		}

		result = append(result, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
