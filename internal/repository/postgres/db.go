package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	url  string
}

func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("empty database URL")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	// Quick connection timeout for resilience
	connectCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	log.Printf("[DB] Connected to PostgreSQL 16 cluster at %s (max_conns: 25)", config.ConnConfig.Host)
	return &DB{pool: pool, url: databaseURL}, nil
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) Ping(ctx context.Context) error {
	if db.pool == nil {
		return fmt.Errorf("postgres pool is nil")
	}
	return db.pool.Ping(ctx)
}

// ApplyMigrations executes the SQL migration script from the given path
func (db *DB) ApplyMigrations(ctx context.Context, migrationPath string) error {
	data, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
	}

	_, err = db.pool.Exec(ctx, string(data))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Printf("[DB] Successfully applied schema migration: %s", migrationPath)
	return nil
}

