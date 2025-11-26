package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	connStr := "postgres://postgres:password@localhost:5433/campaign_db?sslmode=disable"

	// Run migrations first (using database/sql + goose)
	if err := runMigrations(connStr); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Then connect with high-performance pgxpool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Parse config failed:", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL + migrations applied")
	return pool
}
