package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"note-thing/backend/internal/db"

	"github.com/joho/godotenv"
)

const defaultPort = "18611"

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database, err := db.Open()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           routes(database),
	}

	shutdownErrors := make(chan error, 1)
	go func() {
		shutdownErrors <- server.ListenAndServe()
	}()

	log.Printf("backend listening on http://localhost:%s", port)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-shutdownErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	case <-signalChannel:
		log.Printf("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Printf("server stopped")
}

type note struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func routes(database *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/notes", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rows, err := database.QueryContext(
			request.Context(),
			`SELECT id, title, body, created_at FROM notes ORDER BY created_at DESC`,
		)
		if err != nil {
			http.Error(writer, "failed to query notes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		notes := make([]note, 0)
		for rows.Next() {
			var row note
			if err := rows.Scan(&row.Id, &row.Title, &row.Body, &row.CreatedAt); err != nil {
				http.Error(writer, "failed to read notes", http.StatusInternalServerError)
				return
			}
			notes = append(notes, row)
		}

		if err := rows.Err(); err != nil {
			http.Error(writer, "failed while reading notes", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(notes)
	})

	return mux
}
