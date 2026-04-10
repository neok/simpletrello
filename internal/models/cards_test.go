package models_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/neok/simpletrello/internal/models"
)

func setupTabWithCards(t *testing.T) (*models.TabModel, *models.CardModel, models.Tab) {
	t.Helper()
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}
	tab, err := tabs.Insert("My Tab")
	if err != nil {
		t.Fatal(err)
	}
	return tabs, cards, tab
}

func TestCardModel_GetByTab_Empty(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	result, err := cards.GetByTab(tab.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("got %d cards, want 0", len(result))
	}
}

func TestCardModel_Insert(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	card, err := cards.Insert(tab.ID, "Fix bug", "details here")
	if err != nil {
		t.Fatal(err)
	}

	if card.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if card.TabID != tab.ID {
		t.Errorf("got tab_id %d, want %d", card.TabID, tab.ID)
	}
	if card.Title != "Fix bug" {
		t.Errorf("got title %q, want %q", card.Title, "Fix bug")
	}
	if card.Description != "details here" {
		t.Errorf("got description %q, want %q", card.Description, "details here")
	}
	if card.Position != 0 {
		t.Errorf("got position %d, want 0", card.Position)
	}
}

func TestCardModel_Insert_PositionsIncrement(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	for i, title := range []string{"First", "Second", "Third"} {
		card, err := cards.Insert(tab.ID, title, "")
		if err != nil {
			t.Fatal(err)
		}
		if card.Position != i {
			t.Errorf("card %q: got position %d, want %d", title, card.Position, i)
		}
	}
}

func TestCardModel_GetByTab_OnlyReturnsOwnCards(t *testing.T) {
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}

	tabA, _ := tabs.Insert("Tab A")
	tabB, _ := tabs.Insert("Tab B")

	cards.Insert(tabA.ID, "Card A1", "")
	cards.Insert(tabA.ID, "Card A2", "")
	cards.Insert(tabB.ID, "Card B1", "")

	result, err := cards.GetByTab(tabA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("got %d cards for tab A, want 2", len(result))
	}
	for _, c := range result {
		if c.TabID != tabA.ID {
			t.Errorf("card %d has tab_id %d, want %d", c.ID, c.TabID, tabA.ID)
		}
	}
}

func TestCardModel_GetByTab_OrderedByPosition(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	for _, title := range []string{"X", "Y", "Z"} {
		cards.Insert(tab.ID, title, "")
	}

	result, err := cards.GetByTab(tab.ID)
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"X", "Y", "Z"} {
		if result[i].Title != want {
			t.Errorf("result[%d]: got %q, want %q", i, result[i].Title, want)
		}
	}
}

func TestCardModel_Update_Title(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	card, _ := cards.Insert(tab.ID, "Old", "")
	newTitle := "New Title"

	updated, err := cards.Update(card.ID, models.CardUpdate{Title: &newTitle})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "New Title" {
		t.Errorf("got title %q, want %q", updated.Title, "New Title")
	}
	if updated.Description != card.Description {
		t.Error("description should be unchanged")
	}
}

func TestCardModel_Update_Description(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	card, _ := cards.Insert(tab.ID, "Title", "old desc")
	newDesc := "updated desc"

	updated, err := cards.Update(card.ID, models.CardUpdate{Description: &newDesc})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Description != "updated desc" {
		t.Errorf("got description %q, want %q", updated.Description, "updated desc")
	}
	if updated.Title != card.Title {
		t.Error("title should be unchanged")
	}
}

func TestCardModel_Update_MoveToTab(t *testing.T) {
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}

	tabA, _ := tabs.Insert("Tab A")
	tabB, _ := tabs.Insert("Tab B")
	card, _ := cards.Insert(tabA.ID, "Traveler", "")

	updated, err := cards.Update(card.ID, models.CardUpdate{TabID: &tabB.ID})
	if err != nil {
		t.Fatal(err)
	}
	if updated.TabID != tabB.ID {
		t.Errorf("got tab_id %d, want %d", updated.TabID, tabB.ID)
	}

	inA, _ := cards.GetByTab(tabA.ID)
	inB, _ := cards.GetByTab(tabB.ID)
	if len(inA) != 0 {
		t.Errorf("card still in tab A: got %d cards", len(inA))
	}
	if len(inB) != 1 {
		t.Errorf("card not in tab B: got %d cards", len(inB))
	}
}

func TestCardModel_Update_MoveAppendsToEnd(t *testing.T) {
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}

	tabA, _ := tabs.Insert("A")
	tabB, _ := tabs.Insert("B")
	cards.Insert(tabB.ID, "First in B", "")
	cards.Insert(tabB.ID, "Second in B", "")

	card, _ := cards.Insert(tabA.ID, "Moving", "")
	updated, err := cards.Update(card.ID, models.CardUpdate{TabID: &tabB.ID})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Position != 2 {
		t.Errorf("moved card got position %d, want 2 (appended to end)", updated.Position)
	}
}

func TestCardModel_Update_AllFields(t *testing.T) {
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}

	tabA, _ := tabs.Insert("A")
	tabB, _ := tabs.Insert("B")
	card, _ := cards.Insert(tabA.ID, "Original", "original desc")

	newTitle := "Updated"
	newDesc := "updated desc"
	updated, err := cards.Update(card.ID, models.CardUpdate{
		Title:       &newTitle,
		Description: &newDesc,
		TabID:       &tabB.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Updated" {
		t.Errorf("title: got %q", updated.Title)
	}
	if updated.Description != "updated desc" {
		t.Errorf("description: got %q", updated.Description)
	}
	if updated.TabID != tabB.ID {
		t.Errorf("tab_id: got %d", updated.TabID)
	}
}

func TestCardModel_Update_NotFound(t *testing.T) {
	_, cards, _ := setupTabWithCards(t)

	title := "Ghost"
	_, err := cards.Update(999, models.CardUpdate{Title: &title})
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("got %v, want sql.ErrNoRows", err)
	}
}

func TestCardModel_Delete(t *testing.T) {
	_, cards, tab := setupTabWithCards(t)

	card, _ := cards.Insert(tab.ID, "Temp", "")
	if err := cards.Delete(card.ID); err != nil {
		t.Fatal(err)
	}

	result, _ := cards.GetByTab(tab.ID)
	if len(result) != 0 {
		t.Errorf("got %d cards after delete, want 0", len(result))
	}
}

func TestCardModel_Delete_NonExistent(t *testing.T) {
	_, cards, _ := setupTabWithCards(t)

	if err := cards.Delete(999); err != nil {
		t.Errorf("deleting non-existent card should not error, got: %v", err)
	}
}
