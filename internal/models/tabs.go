package models

import (
	"database/sql"
	"time"
)

type Tab struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	Cards     []Card    `json:"cards"`
}

type TabModel struct {
	DB *sql.DB
}

func (m *TabModel) GetAll() ([]Tab, error) {
	tabs := []Tab{}
	rows, err := m.DB.Query(`SELECT id, name, position, created_at FROM tabs ORDER BY position, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Tab
		if err := rows.Scan(&t.ID, &t.Name, &t.Position, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Cards = []Card{}
		tabs = append(tabs, t)
	}
	return tabs, rows.Err()
}

func (m *TabModel) Insert(name string) (Tab, error) {
	var maxPos int
	m.DB.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM tabs`).Scan(&maxPos)

	res, err := m.DB.Exec(`INSERT INTO tabs (name, position) VALUES (?, ?)`, name, maxPos+1)
	if err != nil {
		return Tab{}, err
	}
	id, _ := res.LastInsertId()

	var t Tab
	err = m.DB.QueryRow(`SELECT id, name, position, created_at FROM tabs WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.Position, &t.CreatedAt)
	t.Cards = []Card{}
	return t, err
}

func (m *TabModel) Update(id int64, name string) (Tab, error) {
	_, err := m.DB.Exec(`UPDATE tabs SET name = ? WHERE id = ?`, name, id)
	if err != nil {
		return Tab{}, err
	}
	var t Tab
	err = m.DB.QueryRow(`SELECT id, name, position, created_at FROM tabs WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.Position, &t.CreatedAt)
	t.Cards = []Card{}
	return t, err
}

func (m *TabModel) Delete(id int64) error {
	_, err := m.DB.Exec(`DELETE FROM tabs WHERE id = ?`, id)
	return err
}
