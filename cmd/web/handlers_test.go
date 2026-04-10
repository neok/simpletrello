package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest(t *testing.T, app *application, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *strings.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	} else {
		reqBody = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	app.routes().ServeHTTP(w, req)
	return w
}

func decodeData(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func TestIndex(t *testing.T) {
	app := newTestApp(t)

	w := doRequest(t, app, http.MethodGet, "/", "")

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type: got %q, want text/html", ct)
	}
	if !strings.Contains(w.Body.String(), "__INITIAL_DATA__") {
		t.Error("response missing __INITIAL_DATA__")
	}
}

func TestIndex_NotFound(t *testing.T) {
	app := newTestApp(t)

	w := doRequest(t, app, http.MethodGet, "/no-such-page", "")

	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestListTabs(t *testing.T) {
	app := newTestApp(t)

	t.Run("empty board", func(t *testing.T) {
		w := doRequest(t, app, http.MethodGet, "/api/v1/tabs", "")
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		resp := decodeData(t, w)
		tabs := resp["data"].([]any)
		if len(tabs) != 0 {
			t.Errorf("got %d tabs, want 0", len(tabs))
		}
	})

	t.Run("with tabs and cards", func(t *testing.T) {
		app2 := newTestApp(t)
		tabID := mustInsertTab(t, app2.models.Tabs.DB, "To Do")
		mustInsertCard(t, app2.models.Cards.DB, tabID, "Task 1")
		mustInsertCard(t, app2.models.Cards.DB, tabID, "Task 2")

		w := doRequest(t, app2, http.MethodGet, "/api/v1/tabs", "")
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
		resp := decodeData(t, w)
		tabs := resp["data"].([]any)
		if len(tabs) != 1 {
			t.Fatalf("got %d tabs, want 1", len(tabs))
		}
		tab := tabs[0].(map[string]any)
		cards := tab["cards"].([]any)
		if len(cards) != 2 {
			t.Errorf("got %d cards, want 2", len(cards))
		}
	})
}

func TestCreateTab(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"valid", `{"name":"Backlog"}`, http.StatusCreated},
		{"empty name", `{"name":""}`, http.StatusBadRequest},
		{"missing name field", `{}`, http.StatusBadRequest},
		{"invalid json", `{bad}`, http.StatusBadRequest},
		{"unknown field", `{"name":"x","extra":1}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			w := doRequest(t, app, http.MethodPost, "/api/v1/tabs", tt.body)
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestCreateTab_ReturnsTab(t *testing.T) {
	app := newTestApp(t)
	w := doRequest(t, app, http.MethodPost, "/api/v1/tabs", `{"name":"Sprint"}`)

	resp := decodeData(t, w)
	tab := resp["data"].(map[string]any)
	if tab["name"] != "Sprint" {
		t.Errorf("got name %q, want Sprint", tab["name"])
	}
	if tab["id"] == nil {
		t.Error("expected id in response")
	}
}

func TestUpdateTab(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		body   string
		status int
	}{
		{"valid rename", "1", `{"name":"Renamed"}`, http.StatusOK},
		{"empty name", "1", `{"name":""}`, http.StatusBadRequest},
		{"invalid id", "abc", `{"name":"X"}`, http.StatusBadRequest},
		{"invalid json", "1", `{bad}`, http.StatusBadRequest},
		{"unknown field", "1", `{"name":"X","extra":true}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			mustInsertTab(t, app.models.Tabs.DB, "Original")
			w := doRequest(t, app, http.MethodPatch, "/api/v1/tabs/"+tt.id, tt.body)
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestUpdateTab_NotFound(t *testing.T) {
	app := newTestApp(t)
	w := doRequest(t, app, http.MethodPatch, "/api/v1/tabs/999", `{"name":"Ghost"}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteTab(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		status int
	}{
		{"valid", "1", http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest},
		{"non-existent id", "999", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			mustInsertTab(t, app.models.Tabs.DB, "Temp")
			w := doRequest(t, app, http.MethodDelete, "/api/v1/tabs/"+tt.id, "")
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestDeleteTab_RemovesFromList(t *testing.T) {
	app := newTestApp(t)
	mustInsertTab(t, app.models.Tabs.DB, "Ephemeral")

	doRequest(t, app, http.MethodDelete, "/api/v1/tabs/1", "")

	w := doRequest(t, app, http.MethodGet, "/api/v1/tabs", "")
	resp := decodeData(t, w)
	tabs := resp["data"].([]any)
	if len(tabs) != 0 {
		t.Errorf("got %d tabs after delete, want 0", len(tabs))
	}
}

func TestCreateCard(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"valid", `{"tab_id":1,"title":"New card"}`, http.StatusCreated},
		{"with description", `{"tab_id":1,"title":"Card","description":"details"}`, http.StatusCreated},
		{"missing title", `{"tab_id":1}`, http.StatusBadRequest},
		{"empty title", `{"tab_id":1,"title":""}`, http.StatusBadRequest},
		{"missing tab_id", `{"title":"Card"}`, http.StatusBadRequest},
		{"zero tab_id", `{"tab_id":0,"title":"Card"}`, http.StatusBadRequest},
		{"invalid json", `{bad}`, http.StatusBadRequest},
		{"unknown field", `{"tab_id":1,"title":"X","extra":1}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			mustInsertTab(t, app.models.Tabs.DB, "My Tab")
			w := doRequest(t, app, http.MethodPost, "/api/v1/cards", tt.body)
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestCreateCard_ReturnsCard(t *testing.T) {
	app := newTestApp(t)
	mustInsertTab(t, app.models.Tabs.DB, "Tab")

	w := doRequest(t, app, http.MethodPost, "/api/v1/cards", `{"tab_id":1,"title":"My card","description":"desc"}`)

	resp := decodeData(t, w)
	card := resp["data"].(map[string]any)
	if card["title"] != "My card" {
		t.Errorf("got title %q, want My card", card["title"])
	}
	if card["description"] != "desc" {
		t.Errorf("got description %q, want desc", card["description"])
	}
}

func TestUpdateCard(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		body   string
		status int
	}{
		{"update title", "1", `{"title":"Updated"}`, http.StatusOK},
		{"update description", "1", `{"description":"new desc"}`, http.StatusOK},
		{"update both", "1", `{"title":"T","description":"D"}`, http.StatusOK},
		{"empty body", "1", `{}`, http.StatusOK},
		{"invalid id", "abc", `{"title":"X"}`, http.StatusBadRequest},
		{"invalid json", "1", `{bad}`, http.StatusBadRequest},
		{"unknown field", "1", `{"title":"X","extra":1}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			tabID := mustInsertTab(t, app.models.Tabs.DB, "Tab")
			mustInsertCard(t, app.models.Cards.DB, tabID, "Original")
			w := doRequest(t, app, http.MethodPatch, "/api/v1/cards/"+tt.id, tt.body)
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestUpdateCard_MoveToTab(t *testing.T) {
	app := newTestApp(t)
	tabA := mustInsertTab(t, app.models.Tabs.DB, "A")
	tabB := mustInsertTab(t, app.models.Tabs.DB, "B")
	mustInsertCard(t, app.models.Cards.DB, tabA, "Mover")

	body := fmt.Sprintf(`{"tab_id":%d}`, tabB)
	w := doRequest(t, app, http.MethodPatch, "/api/v1/cards/1", body)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	resp := decodeData(t, w)
	card := resp["data"].(map[string]any)
	if int64(card["tab_id"].(float64)) != tabB {
		t.Errorf("got tab_id %v, want %d", card["tab_id"], tabB)
	}
}

func TestUpdateCard_NotFound(t *testing.T) {
	app := newTestApp(t)
	w := doRequest(t, app, http.MethodPatch, "/api/v1/cards/999", `{"title":"Ghost"}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		status int
	}{
		{"valid", "1", http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest},
		{"non-existent", "999", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(t)
			tabID := mustInsertTab(t, app.models.Tabs.DB, "Tab")
			mustInsertCard(t, app.models.Cards.DB, tabID, "Card")
			w := doRequest(t, app, http.MethodDelete, "/api/v1/cards/"+tt.id, "")
			if w.Code != tt.status {
				t.Errorf("got %d, want %d — body: %s", w.Code, tt.status, w.Body.String())
			}
		})
	}
}

func TestDeleteCard_RemovesFromTab(t *testing.T) {
	app := newTestApp(t)
	tabID := mustInsertTab(t, app.models.Tabs.DB, "Tab")
	mustInsertCard(t, app.models.Cards.DB, tabID, "Gone")

	doRequest(t, app, http.MethodDelete, "/api/v1/cards/1", "")

	w := doRequest(t, app, http.MethodGet, "/api/v1/tabs", "")
	resp := decodeData(t, w)
	tabs := resp["data"].([]any)
	tab := tabs[0].(map[string]any)
	cards := tab["cards"].([]any)
	if len(cards) != 0 {
		t.Errorf("got %d cards after delete, want 0", len(cards))
	}
}
