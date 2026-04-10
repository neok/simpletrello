package models_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/neok/simpletrello/internal/models"
)

func TestTabModel_GetAll_Empty(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}
	tabs, err := m.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(tabs) != 0 {
		t.Errorf("got %d tabs, want 0", len(tabs))
	}
}

func TestTabModel_Insert(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	tab, err := m.Insert("To Do")
	if err != nil {
		t.Fatal(err)
	}

	if tab.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if tab.Name != "To Do" {
		t.Errorf("got name %q, want %q", tab.Name, "To Do")
	}
	if tab.Position != 0 {
		t.Errorf("got position %d, want 0", tab.Position)
	}
	if tab.Cards == nil {
		t.Error("Cards should be initialized, not nil")
	}
}

func TestTabModel_Insert_PositionsIncrement(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	names := []string{"Backlog", "In Progress", "Done"}
	for i, name := range names {
		tab, err := m.Insert(name)
		if err != nil {
			t.Fatal(err)
		}
		if tab.Position != i {
			t.Errorf("tab %q: got position %d, want %d", name, tab.Position, i)
		}
	}
}

func TestTabModel_GetAll_OrderedByPosition(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	for _, name := range []string{"A", "B", "C"} {
		if _, err := m.Insert(name); err != nil {
			t.Fatal(err)
		}
	}

	tabs, err := m.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(tabs) != 3 {
		t.Fatalf("got %d tabs, want 3", len(tabs))
	}
	for i, want := range []string{"A", "B", "C"} {
		if tabs[i].Name != want {
			t.Errorf("tabs[%d]: got %q, want %q", i, tabs[i].Name, want)
		}
	}
}

func TestTabModel_GetAll_CardsInitializedEmpty(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}
	m.Insert("Solo")

	tabs, err := m.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if tabs[0].Cards == nil {
		t.Error("Cards should be an empty slice, not nil")
	}
}

func TestTabModel_Update(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	created, _ := m.Insert("Old Name")
	updated, err := m.Update(created.ID, "New Name")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "New Name" {
		t.Errorf("got %q, want %q", updated.Name, "New Name")
	}
	if updated.ID != created.ID {
		t.Errorf("ID changed: got %d, want %d", updated.ID, created.ID)
	}
}

func TestTabModel_Update_NotFound(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	_, err := m.Update(999, "Ghost")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("got %v, want sql.ErrNoRows", err)
	}
}

func TestTabModel_Delete(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	tab, _ := m.Insert("Temp")
	if err := m.Delete(tab.ID); err != nil {
		t.Fatal(err)
	}

	tabs, _ := m.GetAll()
	if len(tabs) != 0 {
		t.Errorf("got %d tabs after delete, want 0", len(tabs))
	}
}

func TestTabModel_Delete_NonExistent(t *testing.T) {
	m := &models.TabModel{DB: newTestDB(t)}

	if err := m.Delete(999); err != nil {
		t.Errorf("deleting non-existent tab should not error, got: %v", err)
	}
}

func TestTabModel_Delete_CascadesCards(t *testing.T) {
	db := newTestDB(t)
	tabs := &models.TabModel{DB: db}
	cards := &models.CardModel{DB: db}

	tab, _ := tabs.Insert("Work")
	cards.Insert(tab.ID, "Task 1", "")
	cards.Insert(tab.ID, "Task 2", "")

	if err := tabs.Delete(tab.ID); err != nil {
		t.Fatal(err)
	}

	remaining, err := cards.GetByTab(tab.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Errorf("got %d cards after tab delete, want 0 (cascade failed)", len(remaining))
	}
}
