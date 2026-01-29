package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	// Crear carpeta de datos
	_ = os.MkdirAll("data", 0755)

	dbPath := filepath.Join("data", "minicrm.db")
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	DB = database
	return createTables()
}

func createTables() error {
	// Tabla clientes
	clientsQuery := `
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
	if _, err := DB.Exec(clientsQuery); err != nil {
		return err
	}

	// Tabla tareas
	tasksQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		status TEXT,
		priority TEXT,
		owner TEXT,
		progress INTEGER,
		due_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(client_id) REFERENCES clients(id)
	);
	`
	if _, err := DB.Exec(tasksQuery); err != nil {
		return err
	}

	return nil
}
