package main

import (
	"flag"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/neok/simpletrello/internal/models"
	"github.com/neok/simpletrello/ui"
)

type config struct {
	addr         string
	dsn          string
	migrationSQL string
}

type application struct {
	cfg       config
	logger    *slog.Logger
	templates *template.Template
	static    fs.FS
	models    struct {
		Tabs  *models.TabModel
		Cards *models.CardModel
	}
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP listen address")
	flag.StringVar(&cfg.dsn, "dsn", "./data/trello.db", "SQLite DSN")
	flag.StringVar(&cfg.migrationSQL, "migrations", "./migrations/001_initial.sql", "SQL migration file")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tmpl, err := template.ParseFS(ui.Files, "html/*.html")
	if err != nil {
		logger.Error("cannot parse templates", "err", err)
		os.Exit(1)
	}

	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		logger.Error("cannot sub static fs", "err", err)
		os.Exit(1)
	}

	db, err := models.Open(cfg.dsn)
	if err != nil {
		logger.Error("cannot open db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := models.Migrate(db, cfg.migrationSQL); err != nil {
		logger.Error("migration failed", "err", err)
		os.Exit(1)
	}

	app := &application{cfg: cfg, logger: logger, templates: tmpl, static: staticFS}
	app.models.Tabs = &models.TabModel{DB: db}
	app.models.Cards = &models.CardModel{DB: db}

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", "addr", cfg.addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
