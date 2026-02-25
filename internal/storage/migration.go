package storage

type Migration struct {
	Name  string
	Query string
}

var migrations = []Migration{
	{
		Name: "001_create_tasks",
		Query: `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			spent_seconds INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		);
		`,
	},
}