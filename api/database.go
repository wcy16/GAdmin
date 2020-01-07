package api

import (
	"database/sql"
	"fmt"
	"gadmin/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DBConnect() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@/%s",
		config.SETTING.Database.Username,
		config.SETTING.Database.Password,
		config.SETTING.Database.DBname)

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
}

func GetTables() []string {
	rows, err := db.Query("SHOW TABLES")

	if err != nil {
		panic(err)
	}

	var table string
	var tables []string

	for rows.Next() {
		_ = rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables
}

func GetCount(name string) int {
	script := fmt.Sprintf("SELECT COUNT(*) FROM %s", name)
	rows, err := db.Query(script)
	if err != nil {
		panic(err)
	}
	rows.Next()
	var count int
	rows.Scan(&count)
	return count
}

func GetCols(name string) []string {
	script := fmt.Sprintf("SELECT * FROM %s LIMIT 0", name)
	rows, err := db.Query(script)
	if err != nil {
		panic(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	return cols
}

func GetRows(name string, limit, offset uint) (cols []string, results [][]string) {
	script := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", name, limit, offset)
	rows, err := db.Query(script)
	if err != nil {
		panic(err)
	}

	cols, err = rows.Columns()
	if err != nil {
		panic(err)
	}

	results = make([][]string, 0, limit)

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				result[i] = string(raw)
			}
		}

		results = append(results, result)
	}
	return
}
