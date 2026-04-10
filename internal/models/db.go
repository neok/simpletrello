package models

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB, path string) error {
	script, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(script))
	return err
}
