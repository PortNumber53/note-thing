package main

import (
	"flag"
	"fmt"
	"os"

	"note-thing/backend/internal/migrations"
)

func main() {
	var direction string
	var steps int

	flag.StringVar(&direction, "direction", "up", "migration direction: up or down")
	flag.IntVar(&steps, "steps", 0, "number of migration steps; 0 means all")
	flag.Parse()

	options := migrations.RunOptions{
		Direction: migrations.Direction(direction),
		Steps:     steps,
	}

	if err := migrations.Run(options); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("migrations applied successfully")
}
