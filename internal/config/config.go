package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	Environment    string
	RiotAPIKey     string
	GovernorMode   string // DEV_PROTOTYPE_MODE vs PRODUCTION_MODE
	DatabaseURL    string
	RedisURL       string
	OpenSearchURL  string
	JWTSecret      string
	CurrentPatch   string
	DDragonVersion string
	// Comma-separated CORS origins. Empty = same-origin only (safest).
	AllowedOrigins string
	// Seed admin credentials (used when no admin exists yet).
	AdminEmail    string
	AdminPassword string
	// Google OAuth2 (empty = social login disabled).
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	// Sentry error tracking (empty = disabled).
	SentryDSN string
	// Public SPA origin for split deploys (empty = same-origin redirects).
	FrontendBaseURL string
}

func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if os.Getenv(k) == "" {
				_ = os.Setenv(k, v)
			}
		}
	}
}

func Load() *Config {
	loadDotEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	riotKey := os.Getenv("RIOT_API_KEY")

	govMode := os.Getenv("GOVERNOR_MODE")
	if govMode == "" {
		if riotKey == "" || len(riotKey) < 20 {
			govMode = "DEV_PROTOTYPE_MODE"
		} else {
			govMode = "DEV_PROTOTYPE_MODE" // Default development key is 20req/1s
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/riftstats?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "riftstats-super-secure-production-jwt-signing-secret-2026"
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@riftstats.bn.io"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" && env != "production" {
		adminPassword = "securepassword"
	}

	return &Config{
		Port:               port,
		Environment:        env,
		RiotAPIKey:         riotKey,
		GovernorMode:       govMode,
		DatabaseURL:        dbURL,
		RedisURL:           redisURL,
		OpenSearchURL:      os.Getenv("OPENSEARCH_URL"),
		JWTSecret:          jwtSecret,
		CurrentPatch:       "16.21",
		DDragonVersion:     "16.21.1",
		AllowedOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		AdminEmail:         adminEmail,
		AdminPassword:      adminPassword,
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		SentryDSN:          os.Getenv("SENTRY_DSN"),
		FrontendBaseURL:    strings.TrimRight(os.Getenv("FRONTEND_BASE_URL"), "/"),
	}
}

func (c *Config) PortInt() int {
	p, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return p
}
