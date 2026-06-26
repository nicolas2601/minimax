package db

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nicolas/finanzas/backend/internal/config"
)

// Connect opens a GORM connection to Postgres using cfg.DatabaseURL,
// performs a ping to verify reachability, and returns a ready-to-use
// *gorm.DB. The caller owns the connection lifecycle and is responsible
// for closing it via the underlying *sql.DB when shutting down.
//
// All failure modes (empty URL, dial error, ping error) are returned
// to the caller. This package never calls log.Fatal — main.go owns the
// shutdown decision so tests and alternative entry points stay in control.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("db: config is nil")
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("db: DATABASE_URL is empty")
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db: open postgres: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("db: get underlying *sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	log.Println("Database connection established")
	return gormDB, nil
}