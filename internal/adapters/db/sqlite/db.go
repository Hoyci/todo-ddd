package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/app.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	schemas := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			priority INTEGER,
			status TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP
		);
		`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
	}
	return nil
}
