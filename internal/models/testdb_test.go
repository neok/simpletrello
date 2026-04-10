package models_test

import (
	"database/sql"
	"testing"

	"github.com/neok/simpletrello/internal/models"
)

const testSchema = `
CREATE TABLE IF NOT EXISTS tabs (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL,
    position   INTEGER  NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS cards (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    tab_id      INTEGER  NOT NULL REFERENCES tabs(id) ON DELETE CASCADE,
    title       TEXT     NOT NULL,
    description TEXT     NOT NULL DEFAULT '',
    position    INTEGER  NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := models.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(testSchema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
