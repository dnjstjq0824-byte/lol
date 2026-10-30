package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository/postgres"
)

func main() {
	migrationFile := flag.String("file", "migrations/0001_initial_schema.up.sql", "Path to migration SQL file")
	flag.Parse()

	log.Println("==================================================================")
	log.Println("BN RiftStats Database Migration Tool (Directive BN-RIFTSTATS-1000-001)")
	log.Println("PostgreSQL 16 Schema Initializer (50 Tables & Indexes)")
	log.Println("==================================================================")

	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Connecting to database: %s\n", cfg.DatabaseURL)
	db, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v\nEnsure PostgreSQL is running (e.g. `docker compose up -d postgres`)", err)
	}
	defer db.Close()

	if _, err := os.Stat(*migrationFile); os.IsNotExist(err) {
		log.Fatalf("Migration file not found at %s", *migrationFile)
	}

	log.Printf("Applying migration schema from %s ...\n", *migrationFile)
	start := time.Now()
	if err := db.ApplyMigrations(ctx, *migrationFile); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Migration completed successfully in %v.\n", time.Since(start))
	fmt.Println("All 50 canonical tables and indexes are ready.")
}
