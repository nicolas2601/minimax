//go:build integration
// +build integration

package testhelpers

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	DB *gorm.DB
	Cleanup func()
}

func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	pgContainer, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		tcpostgres.WithDatabase("finanzas_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			tcwait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err, "failed to start postgres container")

	baseConnStr, err := pgContainer.ConnectionString(ctx)
	require.NoError(t, err)
	connStr := baseConnStr
	if !strings.Contains(connStr, "sslmode=") {
		sep := "?"
		if strings.Contains(connStr, "?") {
			sep = "&"
		}
		connStr = connStr + sep + "sslmode=disable"
	}

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err)

	if err := runMigrations(connStr); err != nil {
		t.Fatalf("migrations failed: %v", err)
	}

	return &TestDB{
		DB: gormDB,
		Cleanup: func() {
			sqlDB, _ := gormDB.DB()
			if sqlDB != nil {
				_ = sqlDB.Close()
			}
			_ = pgContainer.Terminate(ctx)
		},
	}
}

func runMigrations(dbURL string) error {
	migrationsPath, err := filepath.Abs("../../migrations")
	if err != nil {
		return err
	}
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}