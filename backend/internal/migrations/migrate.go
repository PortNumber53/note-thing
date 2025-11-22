package migrations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"note-thing/backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

type RunOptions struct {
	Direction Direction
	Steps     int
}

func Run(options RunOptions) error {
	_ = config.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required (set it in .env or environment)")
	}

	migrationsDirectory := filepath.Join(".", "migrations")
	absoluteMigrationsDirectory, err := filepath.Abs(migrationsDirectory)
	if err != nil {
		return fmt.Errorf("resolve migrations directory: %w", err)
	}

	sourceURL := "file://" + absoluteMigrationsDirectory

	migrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer migrator.Close()

	switch options.Direction {
	case DirectionUp:
		if options.Steps > 0 {
			err = migrator.Steps(options.Steps)
		} else {
			err = migrator.Up()
		}
	case DirectionDown:
		if options.Steps > 0 {
			err = migrator.Steps(-options.Steps)
		} else {
			err = migrator.Down()
		}
	default:
		return fmt.Errorf("unsupported direction %q", options.Direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
