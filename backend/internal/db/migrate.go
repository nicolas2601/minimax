package db

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/nicolas/finanzas/backend/internal/config"
)

func Migrate(cfg *config.Config) error {
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}