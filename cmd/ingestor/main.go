package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository/postgres"
	"bn-riftstats/internal/riotclient"
)

func main() {
	log.Println("BN RiftStats MATCH-V5 Ingestion Worker starting...")
	cfg := config.Load()
	client := riotclient.NewClient(cfg)

	log.Printf("Ingestion Worker running with Governor Mode: %s\n", cfg.GovernorMode)
	log.Printf("Tracking %d regions across AMERICAS, ASIA, EUROPE, SEA\n", len(riotclient.RegionsTable))

	// Best-effort persistence for heartbeats; the worker runs without it.
	ctx := context.Background()
	var pg *postgres.PostgresRepo
	if cfg.DatabaseURL != "" {
		if db, err := postgres.Connect(ctx, cfg.DatabaseURL); err != nil {
			log.Printf("[Ingestor] PostgreSQL standby (%v); heartbeats stay in logs only", err)
		} else {
			defer db.Close()
			pg = postgres.NewPostgresRepo(db)
		}
	}

	// Ingestion heartbeat loop with graceful shutdown.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	runCycle := func() {
		cycleCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
		defer cancel()
		latest, err := riotclient.FetchLatestDDragonVersion(cycleCtx)
		if err != nil {
			log.Printf("[Ingestor] Data Dragon version check failed: %v", err)
			if pg != nil {
				_ = pg.UpsertSourceFreshness(cycleCtx, "ddragon:global", "GLOBAL", "DATA_DRAGON", cfg.DDragonVersion, -1, "DEGRADED")
			}
			return
		}
		health := "HEALTHY"
		if latest != cfg.DDragonVersion {
			log.Printf("[Ingestor] PATCH DRIFT: observed Data Dragon %s, serving %s", latest, cfg.DDragonVersion)
		} else {
			log.Printf("[Ingestor] Heartbeat: Data Dragon %s in sync, governor=%s, pendingRefresh=%d",
				latest, cfg.GovernorMode, client.PendingRefreshJobs())
		}
		if pg != nil {
			_ = pg.UpsertSourceFreshness(cycleCtx, "ddragon:global", "GLOBAL", "DATA_DRAGON", latest, 0, health)
		}
	}

	runCycle()
	for {
		select {
		case <-ticker.C:
			runCycle()
		case <-quit:
			log.Println("[Ingestor] Shutdown signal received, stopping.")
			return
		}
	}
}
