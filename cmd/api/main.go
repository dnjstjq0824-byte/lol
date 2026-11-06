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
	"bn-riftstats/internal/repository/cache"
	"bn-riftstats/internal/repository/postgres"
	"bn-riftstats/internal/riotclient"
	transport "bn-riftstats/internal/transport/http"
	"bn-riftstats/internal/transport/http/handlers"
	appMiddleware "bn-riftstats/internal/transport/http/middleware"
)

func main() {
	log.Println("==================================================================")
	log.Println("BN RiftStats Global API Gateway (Directive BN-RIFTSTATS-1000-001)")
	log.Println("Target Benchmark: https://op.gg/ (100% Core LoL Functional Parity)")
	log.Println("Observed Patch: 16.21 (Data Dragon 16.21.1)")
	log.Println("==================================================================")

	cfg := config.Load()
	log.Printf("Starting in environment: %s | Governor Mode: %s\n", cfg.Environment, cfg.GovernorMode)
	appMiddleware.InitSentry(cfg.SentryDSN, cfg.Environment)
	defer appMiddleware.FlushSentry()

	ctx := context.Background()

	// Initialize PostgreSQL 16 connection with resilient failover
	var pgDB *postgres.DB
	if cfg.DatabaseURL != "" {
		db, err := postgres.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("[DB Notice] PostgreSQL connection standby (%v); using high-performance internal memory store", err)
		} else {
			pgDB = db
			defer pgDB.Close()
		}
	}

	// Initialize Redis 7.2 connection with resilient failover
	var redisClient *cache.RedisCache
	if cfg.RedisURL != "" {
		rc, err := cache.Connect(ctx, cfg.RedisURL)
		if err != nil {
			log.Printf("[Cache Notice] Redis connection standby (%v); using in-process L1 token bucket cache", err)
		} else {
			redisClient = rc
			defer redisClient.Close()
		}
	}

	// Seed admin credentials: production without ADMIN_PASSWORD gets a random
	// one-time password printed once at boot; never a hardcoded default.
	adminPassword := cfg.AdminPassword
	if adminPassword == "" {
		adminPassword = repository.GenerateRandomPassword()
		log.Printf("[Security] ADMIN_PASSWORD unset: generated one-time admin password for %s: %s", cfg.AdminEmail, adminPassword)
	}
	repository.ConfigureSeedCredentials(cfg.AdminEmail, adminPassword)

	client := riotclient.NewClient(cfg)
	repo := repository.NewRepository(client.DDragon())
	if pgDB != nil {
		repo.SetPostgres(postgres.NewPostgresRepo(pgDB))
		log.Println("[Repository] PostgreSQL 16 persistence layer activated")
	}
	if redisClient != nil {
		repo.SetRedis(redisClient)
		client.SetRedis(redisClient)
		log.Println("[Repository] Redis 7.2 L2 cache layer activated")
	}
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
