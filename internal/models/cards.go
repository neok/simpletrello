package models

import (
	"database/sql"
	"time"
)

type Card struct {
	ID          int64     `json:"id"`
	TabID       int64     `json:"tab_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
}

type CardModel struct {
	DB *sql.DB
}

func (m *CardModel) GetByTab(tabID int64) ([]Card, error) {
	rows, err := m.DB.Query(
		`SELECT id, tab_id, title, description, position, created_at FROM cards WHERE tab_id = ? ORDER BY position, id`,
		tabID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.ID, &c.TabID, &c.Title, &c.Description, &c.Position, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (m *CardModel) Insert(tabID int64, title, description string) (Card, error) {
	var maxPos int
	m.DB.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM cards WHERE tab_id = ?`, tabID).Scan(&maxPos)

	res, err := m.DB.Exec(
		`INSERT INTO cards (tab_id, title, description, position) VALUES (?, ?, ?, ?)`,
		tabID, title, description, maxPos+1,
	)
	if err != nil {
		return Card{}, err
	}
	id, _ := res.LastInsertId()

	var c Card
	err = m.DB.QueryRow(
		`SELECT id, tab_id, title, description, position, created_at FROM cards WHERE id = ?`, id,
	).Scan(&c.ID, &c.TabID, &c.Title, &c.Description, &c.Position, &c.CreatedAt)
	return c, err
}

type CardUpdate struct {
	Title       *string
	Description *string
	TabID       *int64
}

func (m *CardModel) Update(id int64, u CardUpdate) (Card, error) {
	if u.TabID != nil {
		var maxPos int
		m.DB.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM cards WHERE tab_id = ?`, *u.TabID).Scan(&maxPos)
		_, err := m.DB.Exec(`UPDATE cards SET tab_id = ?, position = ? WHERE id = ?`, *u.TabID, maxPos+1, id)
		if err != nil {
			return Card{}, err
		}
	}
	if u.Title != nil || u.Description != nil {
		_, err := m.DB.Exec(
			`UPDATE cards SET
				title       = COALESCE(?, title),
				description = COALESCE(?, description)
			WHERE id = ?`,
			u.Title, u.Description, id,
		)
		if err != nil {
			return Card{}, err
		}
	}

	var c Card
	err := m.DB.QueryRow(
		`SELECT id, tab_id, title, description, position, created_at FROM cards WHERE id = ?`, id,
	).Scan(&c.ID, &c.TabID, &c.Title, &c.Description, &c.Position, &c.CreatedAt)
	return c, err
}

func (m *CardModel) Delete(id int64) error {
	_, err := m.DB.Exec(`DELETE FROM cards WHERE id = ?`, id)
	return err
}
