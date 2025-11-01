package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"peca-back/internal/service"
	"peca-back/internal/store"
	"peca-back/internal/transport"

	_ "github.com/mattn/go-sqlite3"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321") // Astro dev server
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sql.Open("sqlite3", "./entries.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q := `
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image TEXT NOT NULL
		)
	`

	if _, err := db.Exec(q); err != nil {
		log.Fatal(err.Error())
	}

	EntryStore := store.New(db)
	EntryService := service.New(EntryStore)
	EntryHandler := transport.New(EntryService)

	http.HandleFunc("/entradas", EntryHandler.HandleEntries)
	http.HandleFunc("/entradas/", EntryHandler.HandleEntryByID)

	fmt.Println("Servidor ejecutandose en http://localhost:8080")
	fmt.Println("Api endpoints:")
	fmt.Println("GET /entradas - Obtener todas las entradas")
	fmt.Println("POST /entradas - Crear una entrada")
	fmt.Println("GET /entradas/{id} - Obtener una entrada")
	fmt.Println("PUT /entradas/{id} - Actualizar una entrada")
	fmt.Println("DELETE /entradas/{id} - Eliminar una entrada")

	mux := http.NewServeMux()
	mux.HandleFunc("/entradas", EntryHandler.HandleEntries)
	mux.HandleFunc("/entradas/", EntryHandler.HandleEntryByID)

	// Servir archivos estáticos si es necesario
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))

}
