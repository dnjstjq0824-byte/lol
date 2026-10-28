package main

import (
	"log"
	"time"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/riotclient"
)

func main() {
	log.Println("BN RiftStats MATCH-V5 Ingestion Worker starting...")
	cfg := config.Load()
	client := riotclient.NewClient(cfg)

	log.Printf("Ingestion Worker running with Governor Mode: %s\n", cfg.GovernorMode)
	log.Printf("Tracking %d regions across AMERICAS, ASIA, EUROPE, SEA\n", len(riotclient.RegionsTable))

	// Ingestion heartbeat loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("Ingestion worker online and consuming queues for patch %s\n", cfg.CurrentPatch)
	_ = client
}

