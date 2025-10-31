package store

import (
	"database/sql"
	"peca-back/internal/model"
)

type Store interface {
	getAll() ([]*model.Entry, error)
	getById(id int) (*model.Entry, error)
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
