package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/neok/simpletrello/internal/models"
)

type templateData struct {
	InitialData template.JS
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tabs, err := app.models.Tabs.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for i, tab := range tabs {
		cards, err := app.models.Cards.GetByTab(tab.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		tabs[i].Cards = cards
	}

	raw, err := json.Marshal(tabs)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	app.templates.ExecuteTemplate(w, "index.html", templateData{
		InitialData: template.JS(raw),
	})
}

func (app *application) listTabs(w http.ResponseWriter, r *http.Request) {
	tabs, err := app.models.Tabs.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for i, tab := range tabs {
		cards, err := app.models.Cards.GetByTab(tab.ID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		tabs[i].Cards = cards
	}

	app.writeJSON(w, http.StatusOK, tabs)
}

func (app *application) createTab(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Name == "" {
		app.writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tab, err := app.models.Tabs.Insert(input.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, tab)
}

func (app *application) updateTab(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Name == "" {
		app.writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tab, err := app.models.Tabs.Update(id, input.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, tab)
}

func (app *application) deleteTab(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := app.models.Tabs.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (app *application) createCard(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TabID       int64  `json:"tab_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Title == "" {
		app.writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if input.TabID == 0 {
		app.writeError(w, http.StatusBadRequest, "tab_id is required")
		return
	}

	card, err := app.models.Cards.Insert(input.TabID, input.Title, input.Description)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusCreated, card)
}

func (app *application) updateCard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		TabID       *int64  `json:"tab_id"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	card, err := app.models.Cards.Update(id, models.CardUpdate{
		Title:       input.Title,
		Description: input.Description,
		TabID:       input.TabID,
	})
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, card)
}

func (app *application) deleteCard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := app.models.Cards.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
}
