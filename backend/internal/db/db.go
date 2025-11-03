package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func InitDB() *sqlx.DB{
	db, err := sqlx.Open("sqlite", "./sentrinet.db")
	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS scans(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT,
		port INTEGER,
		is_open BOOLEAN,
		duration_ms INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS jobs(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT NOT NULL,
		start_port INTEGER NOT NULL,
		end_port INTEGER NOT NULL,
		interval_seconds INTEGER NOT NULL,
		active INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS notifications(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		scan_id INTEGER,
		type TEXT,
		message TEXT,
		read BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password_hash TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`

	db.MustExec(schema)
	return db
}