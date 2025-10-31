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
	`

	db.MustExec(schema)
	return db
}