package riotclient_test

import (
	"context"
	"testing"
	"time"

	"bn-riftstats/internal/riotclient"
)

func TestRegionRouting16Regions(t *testing.T) {
	requiredRegions := []string{"KR", "NA", "EUW", "EUNE", "JP", "BR", "LAN", "LAS", "OCE", "TR", "RU", "ME", "SG", "PH", "TH", "TW", "VN"}

	for _, code := range requiredRegions {
		reg, err := riotclient.ResolveRegion(code)
		if err != nil {
			t.Fatalf("Failed to resolve required region %s: %v", code, err)
		}
		if reg.UICode != code {
			t.Errorf("Region code mismatch: expected %s, got %s", code, reg.UICode)
		}
		if reg.PlatformHost == "" || reg.RegionalHost == "" {
			t.Errorf("Missing hosts for region %s", code)
		}
	}

	// Case insensitivity
	reg, err := riotclient.ResolveRegion("kr")
	if err != nil || reg.UICode != "KR" {
		t.Fatalf("Expected case-insensitive match for 'kr'")
	}
}

func TestRateLimiterDevPrototypeMode(t *testing.T) {
	limiter := riotclient.NewRateLimiter("DEV_PROTOTYPE_MODE")
	ctx := context.Background()

	// Should allow initial 20 requests
	for i := 0; i < 20; i++ {
		if err := limiter.Acquire(ctx, "KR", riotclient.PriorityP0); err != nil {
			t.Fatalf("Failed acquiring token %d: %v", i, err)
		}
	}

	// Next immediate request should be denied
	err := limiter.Acquire(ctx, "KR", riotclient.PriorityP0)
	if err == nil {
		t.Fatalf("Expected rate limit exhaustion after 20 requests in dev mode")
	}

	// Test 429 Retry-After handling
	limiter.Record429("KR", 2)
	err429 := limiter.Acquire(ctx, "KR", riotclient.PriorityP0)
	if err429 == nil {
		t.Fatalf("Expected 429 block during cooldown period")
	}

	time.Sleep(2100 * time.Millisecond)
	// Now should recover
	if err := limiter.Acquire(ctx, "KR", riotclient.PriorityP0); err != nil {
		t.Fatalf("Expected recovery after 429 cooldown: %v", err)
	}
}
