package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository"
	"bn-riftstats/internal/repository/postgres"
	"bn-riftstats/internal/riotclient"
)

func main() {
	log.Println("BN RiftStats 7D Statistical Aggregation Worker starting...")
	cfg := config.Load()
	client := riotclient.NewClient(cfg)
	repo := repository.NewRepository(client.DDragon())

	interval := 60 * time.Minute
	if raw := os.Getenv("AGGREGATE_INTERVAL_MIN"); raw != "" {
		if minutes, err := strconv.Atoi(raw); err == nil && minutes >= 1 {
			interval = time.Duration(minutes) * time.Minute
		}
	}
	log.Printf("Running Wilson Score Shrinkage Tier Engine for Patch %s (cycle: %v)\n", cfg.CurrentPatch, interval)
	log.Printf("Total Matches in Ingestion Universe: %d\n", repo.TotalMatches())

	ctx := context.Background()
	var pg *postgres.PostgresRepo
	if cfg.DatabaseURL != "" {
		if db, err := postgres.Connect(ctx, cfg.DatabaseURL); err != nil {
			log.Printf("[Aggregator] PostgreSQL standby (%v); job records stay in logs only", err)
		} else {
			defer db.Close()
			pg = postgres.NewPostgresRepo(db)
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	runCycle := func() {
		processed, changed := repo.RecalculateTierGrades()
		log.Printf("[Aggregator] Cycle complete: %d slices re-scored, %d grade changes (patch %s)",
			processed, changed, cfg.CurrentPatch)
		if pg != nil {
			cycleCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			_ = pg.InsertAnalyticsJob(cycleCtx, "TIER_RECALC", cfg.CurrentPatch, "GLOBAL", "COMPLETED", processed, "")
		}
	}

	runCycle()
	for {
		select {
		case <-ticker.C:
			runCycle()
		case <-quit:
			log.Println("[Aggregator] Shutdown signal received, stopping.")
			return
		}
	}
}
