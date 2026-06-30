package db

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

	// When connecting through Supabase's Supavisor pooler (port 6543), prepared
	// statement caching collides with the pooler's own statement cache and
	// surfaces as SQLSTATE 42P05 "prepared statement ... already exists".
	// Two layered fixes:
	//   1. Disable gorm's prepared-statement cache.
	//   2. Tell pgx (via DSN) to use simple protocol — no server-side prepared
	//      statements — which is the Supabase-recommended setting for pooler.
	// For direct connections (port 5432) both options are harmless no-ops.
	dsn := cfg.DatabaseURL
	if !strings.Contains(dsn, "default_query_exec_mode=") {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		dsn = dsn + sep + "default_query_exec_mode=simple_protocol"
	}
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
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