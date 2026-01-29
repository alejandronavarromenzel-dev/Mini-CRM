package db

func createTasksTable() error {
	query := `
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
	_, err := DB.Exec(query)
	return err
}
