package riotclient

import (
	"context"
	"testing"
	"time"

	"bn-riftstats/internal/config"
)

func TestRateLimiter_GovernorModes(t *testing.T) {
	// DEV_PROTOTYPE_MODE check
	devLimiter := NewRateLimiter("DEV_PROTOTYPE_MODE")
	devBucket := devLimiter.getBucket("KR")
	if devBucket.capacity != 20 {
		t.Errorf("expected DEV_PROTOTYPE_MODE capacity 20, got %d", devBucket.capacity)
	}
	if devBucket.refillRate != 20.0 {
		t.Errorf("expected DEV_PROTOTYPE_MODE refillRate 20.0, got %f", devBucket.refillRate)
	}

	// PRODUCTION_MODE check
	prodLimiter := NewRateLimiter("PRODUCTION_MODE")
	prodBucket := prodLimiter.getBucket("KR")
	if prodBucket.capacity != 500 {
		t.Errorf("expected PRODUCTION_MODE capacity 500, got %d", prodBucket.capacity)
	}
	if prodBucket.refillRate != 50.0 {
		t.Errorf("expected PRODUCTION_MODE refillRate 50.0, got %f", prodBucket.refillRate)
	}
}

func TestRateLimiter_PriorityTiers(t *testing.T) {
	ctx := context.Background()
	rl := NewRateLimiter("DEV_PROTOTYPE_MODE")

	// Consume down to 5 tokens (out of 20)
	bucket := rl.getBucket("KR")
	bucket.tokens = 5.0

	// P2 should fail because it requires >= 8 tokens (40% of 20)
	errP2 := rl.Acquire(ctx, "KR", PriorityP2)
	if errP2 == nil {
		t.Errorf("expected PriorityP2 to be denied at 5 tokens, got nil")
	}

	// P1 should succeed because 5 tokens >= 4 tokens (20% of 20)
	errP1 := rl.Acquire(ctx, "KR", PriorityP1)
	if errP1 != nil {
		t.Errorf("expected PriorityP1 to succeed, got %v", errP1)
	}

	// Drain down to 0.5 tokens
	bucket.tokens = 0.5
	// P0 requires min 1.0 token
	errP0 := rl.Acquire(ctx, "KR", PriorityP0)
	if errP0 == nil {
		t.Errorf("expected PriorityP0 to fail with insufficient tokens, got nil")
	}
}

func TestRateLimiter_429RetryAfter(t *testing.T) {
	ctx := context.Background()
	rl := NewRateLimiter("DEV_PROTOTYPE_MODE")

	// Trigger 429 for 2 seconds
	rl.Record429("KR", 2)

	// Acquire should fail with 429 message
	err := rl.Acquire(ctx, "KR", PriorityP0)
	if err == nil {
		t.Fatalf("expected error during 429 cooldown, got nil")
	}

	// Simulate cooldown expiry and token replenishment
	bucket := rl.getBucket("KR")
	bucket.retryAfterUntil = time.Now().Add(-1 * time.Second)
	bucket.tokens = 10.0

	errAfter := rl.Acquire(ctx, "KR", PriorityP0)
	if errAfter != nil {
		t.Errorf("expected Acquire to succeed after cooldown, got %v", errAfter)
	}
}

func TestClient_FallbackWithoutKey(t *testing.T) {
	cfg := &config.Config{
		RiotAPIKey:     "",
		GovernorMode:   "DEV_PROTOTYPE_MODE",
		CurrentPatch:   "16.21",
		DDragonVersion: "16.21.1",
	}
	client := NewClient(cfg)
	if client.HasLiveKey() {
		t.Errorf("expected HasLiveKey to be false with empty key")
	}

	ctx := context.Background()
	acc, err := client.ResolveRiotID(ctx, "KR", "Faker", "KR1")
	if err != nil {
		t.Fatalf("expected fallback resolve to succeed, got %v", err)
	}
	if acc.GameName != "Faker" || acc.TagLine != "KR1" {
		t.Errorf("unexpected account data: %+v", acc)
	}

	profile, err := client.GetSummonerProfile(ctx, "KR", acc.Puuid)
	if err != nil {
		t.Fatalf("expected fallback profile to succeed, got %v", err)
	}
	if profile.SummonerLevel <= 0 {
		t.Errorf("expected positive summoner level, got %d", profile.SummonerLevel)
	}
}

