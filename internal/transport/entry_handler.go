package transport

import (
	"encoding/json"
	"net/http"
	"peca-back/internal/model"
	"peca-back/internal/service"
	"strconv"
	"strings"
)

type EntryHandler struct {
	service *service.Service
}

func New(s *service.Service) *EntryHandler {
	return &EntryHandler{
		service: s,
	}
}

func (h *EntryHandler) HandleEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		entries, err := h.service.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB máximo
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Archivo demasiado grande", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "No se encontró la imagen", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename, err := h.service.SaveImage(file, handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		entry := model.Entry{
			Title:   title,
			Content: content,
			Image:   filename,
		}

		created, err := h.service.Create(entry)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)

	default:
		http.Error(w, "Método no disponible", http.StatusBadRequest)
	}
}

func (h *EntryHandler) HandleEntryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/entradas/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "No se encuentra", http.StatusBadRequest)
	}

	switch r.Method {
	case http.MethodGet:
		entry, err := h.service.GetById(id)
		if err != nil {
			http.Error(w, "No se encuentra", http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entry)
	case http.MethodPut:
		var entry model.Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Input inválido", http.StatusBadRequest)
			return
		}
		updated, err := h.service.Update(id, entry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		if err := h.service.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Método no disponible", http.StatusMethodNotAllowed)
	}
}
