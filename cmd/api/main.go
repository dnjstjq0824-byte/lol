package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository"
	"bn-riftstats/internal/riotclient"
	transport "bn-riftstats/internal/transport/http"
	"bn-riftstats/internal/transport/http/handlers"
)

func main() {
	log.Println("==================================================================")
	log.Println("BN RiftStats Global API Gateway (Directive BN-RIFTSTATS-1000-001)")
	log.Println("Target Benchmark: https://op.gg/ (100% Core LoL Functional Parity)")
	log.Println("Observed Patch: 16.21 (Data Dragon 16.21.1)")
	log.Println("==================================================================")

	cfg := config.Load()
	log.Printf("Starting in environment: %s | Governor Mode: %s\n", cfg.Environment, cfg.GovernorMode)

	client := riotclient.NewClient(cfg)
	repo := repository.NewRepository(client.DDragon())
	handler := handlers.NewHandler(cfg, client, repo)
	router := transport.NewRouter(cfg, handler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("HTTP Server listening on %s (OpenAPI 3.1: http://localhost:%s/api/v1)\n", addr, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped.")
}

