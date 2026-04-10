package main

import (
	"database/sql"
	"html/template"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/neok/simpletrello/internal/models"
	_ "modernc.org/sqlite"
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

const indexTmpl = `<!DOCTYPE html><html><body>
<script>window.__INITIAL_DATA__ = {{.InitialData}};</script>
</body></html>`

func newTestApp(t *testing.T) *application {
	t.Helper()

	db, err := models.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(testSchema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	tmpl := template.Must(template.New("index.html").Parse(indexTmpl))

	app := &application{
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		templates: tmpl,
		static:    os.DirFS("."),
	}
	app.models.Tabs = &models.TabModel{DB: db}
	app.models.Cards = &models.CardModel{DB: db}
	return app
}

func mustInsertTab(t *testing.T, db *sql.DB, name string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO tabs (name, position) VALUES (?, 0)`, name)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}

func mustInsertCard(t *testing.T, db *sql.DB, tabID int64, title string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO cards (tab_id, title, description, position) VALUES (?, ?, '', 0)`, tabID, title)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}
