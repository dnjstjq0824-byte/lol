package main

import (
	"log"
	"time"

	"bn-riftstats/internal/analytics"
	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository"
	"bn-riftstats/internal/riotclient"
)

func main() {
	log.Println("BN RiftStats 7D Statistical Aggregation Worker starting...")
	cfg := config.Load()
	client := riotclient.NewClient(cfg)
	repo := repository.NewRepository(client.DDragon())

	log.Printf("Running Wilson Score Shrinkage Tier Engine for Patch %s\n", cfg.CurrentPatch)

	// Test calculation
	wilsonWR := analytics.WilsonScoreLowerBound(26400, 50000)
	score := analytics.CalculateCompositeScore(wilsonWR, 12.4, 8.5, 50.0, 2.5, 12.0, 3.5, 10.0, 5.0)
	grade := analytics.ClassifyTierGrade(score, 12.4)
	quality := analytics.ClassifySampleQuality(50000)

	log.Printf("Sample Ahri Calculation: Wilson WR=%.2f%%, Score=%.2f, Grade=%s, Quality=%s\n",
		wilsonWR, score, grade, quality)
	log.Printf("Total Matches in Ingestion Universe: %d\n", repo.TotalMatches())

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()
	log.Println("Aggregator worker scheduled (60m cycle).")
}

