package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// SecurityHeaders sets baseline hardening headers for every response.
// The SPA is served same-origin, so a strict policy is safe.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

type ipBucket struct {
	tokens     float64
	lastRefill time.Time
}

// IPRateLimiter is a lightweight in-process token bucket keyed by client IP.
// It protects the API (especially auth + Riot-backed fan-out) from abuse.
type IPRateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*ipBucket
	capacity   float64
	refillRate float64 // tokens per second
}

func NewIPRateLimiter(requestsPerMinute int) *IPRateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 300
	}
	return &IPRateLimiter{
		buckets:    make(map[string]*ipBucket),
		capacity:   float64(requestsPerMinute),
		refillRate: float64(requestsPerMinute) / 60.0,
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if host, _, err := net.SplitHostPort(xff); err == nil {
			return host
		}
		return xff
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

func (l *IPRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[ip]
	if !ok {
		b = &ipBucket{tokens: l.capacity, lastRefill: now}
		l.buckets[ip] = b
	}
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * l.refillRate
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.lastRefill = now
	if b.tokens < 1.0 {
		return false
	}
	b.tokens -= 1.0
	return true
}

// Middleware enforces the limit and answers 429 with RFC 9457 problem JSON.
func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.allow(clientIP(r)) {
			w.Header().Set("Content-Type", "application/problem+json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"type":"https://api.riftstats.bn.io/errors/rate-limited","title":"Rate Limited","status":429,"detail":"Too many requests from this IP, slow down and retry.","retryAfterSeconds":60}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
