package store

import (
	"database/sql"
	"peca-back/internal/model"
)

type Store interface {
	GetAll() ([]*model.Entry, error)
	GetById(id int) (*model.Entry, error)
	Create(entry *model.Entry) (*model.Entry, error)
	Update(id int, entry *model.Entry) (*model.Entry, error)
	Delete(id int) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll() ([]*model.Entry, error) {
	q := `SELECT id, title, content, image FROM entries`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []*model.Entry{}
	for rows.Next() {
		entry := model.Entry{}

		err := rows.Scan(&entry.ID, &entry.Title, &entry.Content, &entry.Image)
		if err != nil {
			return nil, err
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

func (s *store) GetById(id int) (*model.Entry, error) {
	q := `SELECT id, title, content, image FROM entries WHERE id = ?`

	entry := model.Entry{}
	err := s.db.QueryRow(q, id).Scan(&entry.ID, &entry.Title, &entry.Content, &entry.Image)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *store) Create(entry *model.Entry) (*model.Entry, error) {
	q := `INSERT INTO entries (title, content, image) VALUES (?, ?, ?)`

	result, err := s.db.Exec(q, entry.Title, entry.Content, entry.Image)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	entry.ID = int(id)

	return entry, nil

}

func (s *store) Update(id int, entry *model.Entry) (*model.Entry, error) {
	q := "UPDATE entries SET title = ?, content = ?, image = ? WHERE id = ?"

	_, err := s.db.Exec(q, entry.Title, entry.Content, entry.Image, id)
	if err != nil {
		return nil, err
	}

	entry.ID = id

	return entry, nil
}

func (s *store) Delete(id int) error {
	q := "DELETE from entries WHERE id = ?"

	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
