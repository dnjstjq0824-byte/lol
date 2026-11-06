package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry enables error tracking when a DSN is configured. With an empty
// DSN it is a no-op so local, test, and fixture runs stay dependency-free.
func InitSentry(dsn, environment string) {
	if dsn == "" {
		return
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
	}); err != nil {
		log.Printf("[Sentry] init failed: %v", err)
		return
	}
	log.Println("[Sentry] Error tracking enabled")
}

// FlushSentry drains buffered events before shutdown.
func FlushSentry() {
	sentry.Flush(3 * time.Second)
}

// SentryReporter captures panics and reports them to Sentry, then re-panics
// so chi's Recoverer still renders the 500 response. It must be registered
// AFTER Recoverer (inner) to observe the panic first. No-op without a DSN.
func SentryReporter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sentry.CurrentHub().Client() == nil {
			next.ServeHTTP(w, r)
			return
		}
		hub := sentry.GetHubFromContext(r.Context())
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
		}
		hub.Scope().SetRequest(r)
		defer func() {
			if recovered := recover(); recovered != nil {
				hub.Recover(recovered)
				sentry.Flush(2 * time.Second)
				panic(recovered)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
