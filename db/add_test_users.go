package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5436 user=api_tester password=testing dbname=film_api sslmode=disable"
	db, err := openDB(connStr)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	stmt := `
	INSERT INTO users (role, api_key) VALUES 
											  ('admin', 'root'),
	('user', '12345');`

	lastInsertId := 0
	err = db.QueryRow(stmt).Scan(&lastInsertId)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn) // right or not?
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
