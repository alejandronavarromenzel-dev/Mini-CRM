package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	os.MkdirAll("data", 0755)

	dbPath := filepath.Join("data", "minicrm.db")
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	DB = database
	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT,
		owner TEXT,
		tags TEXT,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(query)
	return err
}
