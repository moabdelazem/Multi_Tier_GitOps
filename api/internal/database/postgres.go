package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/moabdelazem/mutlitier_app/internal/config"
)

// DB is a wrapper around sql.DB
type DB struct {
	*sql.DB
}

// NewPostgresConnection creates a new PostgreSQL connection
func NewPostgresConnection(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Verify connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully, let's rock!")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	log.Println("Closing database connection...")
	return db.DB.Close()
}

// HealthCheck checks if database is accessible
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.PingContext(ctx)
}

// GetStats returns database connection pool statistics
func (db *DB) GetStats() sql.DBStats {
	return db.Stats()
}
