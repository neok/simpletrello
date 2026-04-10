package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/tabs", app.listTabs)
	mux.HandleFunc("POST /api/v1/tabs", app.createTab)
	mux.HandleFunc("PATCH /api/v1/tabs/{id}", app.updateTab)
	mux.HandleFunc("DELETE /api/v1/tabs/{id}", app.deleteTab)

	mux.HandleFunc("POST /api/v1/cards", app.createCard)
	mux.HandleFunc("PATCH /api/v1/cards/{id}", app.updateCard)
	mux.HandleFunc("DELETE /api/v1/cards/{id}", app.deleteCard)

	mux.Handle("GET /assets/", http.FileServer(http.FS(app.static)))
	mux.HandleFunc("GET /", app.index)

	return mux
}
