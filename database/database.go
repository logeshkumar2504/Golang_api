package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./cameras.db")
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS cameras (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		port INTEGER NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		stream_url TEXT NOT NULL,
		location TEXT,
		status TEXT DEFAULT 'offline',
		camera_type TEXT,
		resolution TEXT,
		frame_rate INTEGER,
		manufacturer TEXT,
		model TEXT
	)`

	_, err = DB.Exec(query)
	return err
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
