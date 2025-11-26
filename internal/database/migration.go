package database

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// runMigrations reads migrations from the real filesystem (100% reliable)
func runMigrations(dbURL string) error {
	// Find the project root (works from any location)
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(currentFile)), "..")
	migrationDir := filepath.Join(projectRoot, "migrations")

	// Debug: show the path
	log.Printf("Looking for migrations in: %s", migrationDir)

	// Open DB connection for goose
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Set dialect
	goose.SetDialect("postgres")

	// Run migrations from real folder
	if err := goose.Up(db, migrationDir); err != nil {
		return err
	}

	log.Println("Migrations applied successfully!")
	return nil
}
