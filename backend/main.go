package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"note-thing/backend/internal/config"
	"note-thing/backend/internal/db"
	"note-thing/backend/internal/migrations"
	"note-thing/backend/internal/router"
)

const defaultPort = "18611"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		runServe()
	case "migrate":
		runMigrate(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: note-thing <command> [options]

Commands:
  serve                  Start the HTTP server
  migrate up             Apply all pending migrations
  migrate up --steps N   Apply N migrations
  migrate down           Rollback all migrations
  migrate down --steps N Rollback N migrations
`)
}

// serve

func runServe() {
	if err := config.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	database, err := db.Open()
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	handler := router.New(database, jwtSecret)

	server := &http.Server{
		Addr:              "0.0.0.0:" + port,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           handler,
	}

	shutdownErrors := make(chan error, 1)
	go func() {
		shutdownErrors <- server.ListenAndServe()
	}()

	log.Printf("backend listening on http://0.0.0.0:%s", port)

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

// migrate

func runMigrate(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: note-thing migrate <up|down> [--steps N]")
		os.Exit(1)
	}

	direction := args[0]
	if direction != "up" && direction != "down" {
		fmt.Fprintf(os.Stderr, "unknown migration direction: %s\n", direction)
		fmt.Fprintln(os.Stderr, "usage: note-thing migrate <up|down> [--steps N]")
		os.Exit(1)
	}

	steps := 0
	for i := 1; i < len(args); i++ {
		if args[i] == "--steps" {
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--steps requires a value")
				os.Exit(1)
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid --steps value: %s\n", args[i+1])
				os.Exit(1)
			}
			steps = n
			i++
		}
	}

	options := migrations.RunOptions{
		Direction: migrations.Direction(direction),
		Steps:     steps,
	}

	if err := migrations.Run(options); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("migrations applied successfully")
}
