package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"peca-back/internal/model"
	"peca-back/internal/store"
	"time"
)

type Service struct {
	store store.Store
}

func New(s store.Store) *Service {
	return &Service{
		store: s,
	}
}

func (s *Service) GetAll() ([]*model.Entry, error) {
	entries, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *Service) GetById(id int) (*model.Entry, error) {
	return s.store.GetById(id)
}

func (s *Service) Create(entry model.Entry) (*model.Entry, error) {
	if entry.Title == "" {
		return nil, errors.New("Se necesita un título")
	}
	if entry.Content == "" {
		return nil, errors.New("Se necesita un contenido")
	}

	return s.store.Create(&entry)
}

func (s *Service) Update(id int, entry model.Entry) (*model.Entry, error) {
	if entry.Title == "" {
		return nil, errors.New("Se necesita un título")
	}
	if entry.Content == "" {
		return nil, errors.New("Se necesita un contenido")
	}

	return s.store.Update(id, &entry)
}

func (s *Service) Delete(id int) error {
	return s.store.Delete(id)
}

func (s *Service) SaveImage(file multipart.File, original string) (string, error) {
	ext := filepath.Ext(original)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("uploads", filename)
	os.MkdirAll("uploads", os.ModePerm)
	dst, _ := os.Create(path)
	defer dst.Close()
	io.Copy(dst, file)
	return filename, nil
}
