package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type DatabaseManager struct {
	db DB
}

func NewDatabaseManager() (*DatabaseManager, error) {
	dbName := "Dota"
	if os.Getenv("GO_ENV") == "test" {
		dbName = "DotaTest"
	}

	db, err := NewDB("localhost", "5432", "postgres", "1q2ws3edc4r", dbName)
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
	// Составляем запрос с плейсхолдерами
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(params, ", "), tableName)

	var whereClauses []string
	var queryArgs []interface{}
	i := 1
	for key, values := range args {
		placeholders := make([]string, len(values))
		for j, value := range values {
			placeholders[j] = fmt.Sprintf("$%d", i)
			queryArgs = append(queryArgs, value)
			i++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Выполняем подготовленный запрос
	rows, err := dbManager.db.Query(query, queryArgs...)
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

// Модифицированный метод для безопасного добавления строк
func (dbManager *DatabaseManager) AddRows(tableName string, args map[string][]string) error {
	columns := make([]string, 0, len(args))
	for col := range args {
		columns = append(columns, col)
	}

	numRows := len(args[columns[0]])
	for _, values := range args {
		if len(values) != numRows {
			return fmt.Errorf("mismatched number of values in columns")
		}
	}

	// Формируем запрос с плейсхолдерами
	valuePlaceholders := make([]string, numRows)
	var queryArgs []interface{}
	argIndex := 1
	for i := 0; i < numRows; i++ {
		rowPlaceholders := make([]string, len(columns))
		for j, col := range columns {
			rowPlaceholders[j] = fmt.Sprintf("$%d", argIndex)
			queryArgs = append(queryArgs, args[col][i])
			argIndex++
		}
		valuePlaceholders[i] = fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", "))
	}

	query := fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES %s", tableName, strings.Join(columns, ", "), strings.Join(valuePlaceholders, ", "))
	// Выполняем подготовленный запрос
	_, err := dbManager.db.Exec(query, queryArgs...)
	if err != nil {
		return fmt.Errorf("error inserting rows: %v", err)
	}

	return nil
}

func (dbManager *DatabaseManager) DeleteRows(tableName string, args map[string][]string) error {
	// Начинаем с основного запроса
	query := fmt.Sprintf("DELETE FROM %s", tableName)

	// Формируем условия WHERE с использованием плейсхолдеров
	var whereClauses []string
	var queryArgs []interface{}
	i := 1
	for key, values := range args {
		placeholders := make([]string, len(values))
		for j, value := range values {
			placeholders[j] = fmt.Sprintf("$%d", i)
			queryArgs = append(queryArgs, value)
			i++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
	}

	// Добавляем условия WHERE, если они есть
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Выполняем запрос
	_, err := dbManager.db.Exec(query, queryArgs...)
	if err != nil {
		return fmt.Errorf("ошибка удаления строк: %v", err)
	}

	return nil
}
