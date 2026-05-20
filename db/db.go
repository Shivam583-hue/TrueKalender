package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var conn *sql.DB

func Init() error {
	var err error

	conn, err = sql.Open("sqlite3", "./events.db")
	if err != nil {
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER NOT NULL,
		month INTEGER NOT NULL,
		day INTEGER NOT NULL
	);`

	if _, err = conn.Exec(query); err != nil {
		return err
	}

	if _, err = conn.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_date ON tasks(year, month, day)"); err != nil {
		return err
	}

	return nil
}

func Insert(title string, year, month, day int) error {
	if conn == nil {
		return fmt.Errorf("database is not initialized")
	}

	_, err := conn.Exec(
		"INSERT INTO tasks(title, year, month, day) VALUES(?, ?, ?, ?)",
		title,
		year,
		month,
		day,
	)
	return err
}

func FetchByDate(year, month, day int) ([]string, error) {
	if conn == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	rows, err := conn.Query(
		"SELECT title FROM tasks WHERE year = ? AND month = ? AND day = ? ORDER BY id ASC",
		year,
		month,
		day,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return titles, nil
}

func CountByDate(year, month, day int) (int, error) {
	if conn == nil {
		return 0, fmt.Errorf("database is not initialized")
	}

	var count int
	err := conn.QueryRow(
		"SELECT COUNT(*) FROM tasks WHERE year = ? AND month = ? AND day = ?",
		year,
		month,
		day,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func Close() error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}
