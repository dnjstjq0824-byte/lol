package riotclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"bn-riftstats/internal/domain"
)

type PriorityTier int

const (
	PriorityP0 PriorityTier = 0 // Interactive Search / Live Game (65% capacity)
	PriorityP1 PriorityTier = 1 // Refresh / Match History (25% capacity)
	PriorityP2 PriorityTier = 2 // Background Crawler / Aggregator (10% capacity)
)

type RateBucket struct {
	mu           sync.Mutex
	platformID   string
	capacity     int
	tokens       float64
	refillRate   float64 // tokens per second
	lastRefill   time.Time
	appLimitStr  string
	retryAfterUntil time.Time
}

type RateLimiter struct {
	mu      sync.RWMutex
	mode    string
	buckets map[string]*RateBucket
}

func NewRateLimiter(mode string) *RateLimiter {
	rl := &RateLimiter{
		mode:    mode,
		buckets: make(map[string]*RateBucket),
	}
	return rl
}

func (rl *RateLimiter) getBucket(platformID string) *RateBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if b, ok := rl.buckets[platformID]; ok {
		return b
	}

	cap := 20
	refill := 20.0
	limitStr := "20:1,100:120"
	if rl.mode == "PRODUCTION_MODE" {
		cap = 500
		refill = 50.0 // 500 per 10s = 50/s
		limitStr = "500:10,30000:600"
	}

	b := &RateBucket{
		platformID:  platformID,
		capacity:    cap,
		tokens:      float64(cap),
		refillRate:  refill,
		lastRefill:  time.Now(),
		appLimitStr: limitStr,
	}
	rl.buckets[platformID] = b
	return b
}

func (rl *RateLimiter) Acquire(ctx context.Context, platformID string, priority PriorityTier) error {
	bucket := rl.getBucket(platformID)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Check if in 429 Retry-After cooldown
	now := time.Now()
	if now.Before(bucket.retryAfterUntil) {
		waitDur := bucket.retryAfterUntil.Sub(now)
		return fmt.Errorf("rate limited (429 cooldown active for %v)", waitDur.Round(time.Second))
	}

	// Refill tokens
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > float64(bucket.capacity) {
		bucket.tokens = float64(bucket.capacity)
	}
	bucket.lastRefill = now

	// Priority thresholds:
	// P2 requires at least 40% capacity remaining
	// P1 requires at least 20% capacity remaining
	// P0 can use up to 100%
	minTokens := 1.0
	switch priority {
	case PriorityP2:
		minTokens = float64(bucket.capacity) * 0.4
	case PriorityP1:
		minTokens = float64(bucket.capacity) * 0.2
	case PriorityP0:
		minTokens = 1.0
	}

	if bucket.tokens < minTokens {
		return fmt.Errorf("rate limit bucket for %s exhausted for priority %d (available: %.1f, required: %.1f)", platformID, priority, bucket.tokens, minTokens)
	}

	bucket.tokens -= 1.0
	return nil
}

func (rl *RateLimiter) Record429(platformID string, retryAfterSec int) {
	bucket := rl.getBucket(platformID)
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if retryAfterSec <= 0 {
		retryAfterSec = 5
	}
	bucket.retryAfterUntil = time.Now().Add(time.Duration(retryAfterSec) * time.Second)
	bucket.tokens = 0
}

func (rl *RateLimiter) GetStatus() []domain.RateLimitBucketStatus {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	list := make([]domain.RateLimitBucketStatus, 0, len(rl.buckets))
	for pid, b := range rl.buckets {
		b.mu.Lock()
		rem := int(b.tokens)
		status := "NORMAL"
		if time.Now().Before(b.retryAfterUntil) {
			status = "PAUSED_429"
		} else if float64(rem) < float64(b.capacity)*0.2 {
			status = "THROTTLED"
		}
		list = append(list, domain.RateLimitBucketStatus{
			PlatformID: pid,
			AppLimit:   b.appLimitStr,
			AppCount:   b.capacity - rem,
			Capacity:   b.capacity,
			Remaining:  rem,
			P0Reserved: int(float64(b.capacity) * 0.65),
			P1Reserved: int(float64(b.capacity) * 0.25),
			P2Reserved: int(float64(b.capacity) * 0.10),
			Status:     status,
		})
		b.mu.Unlock()
	}
	return list
}

