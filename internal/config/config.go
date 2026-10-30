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

	return &Config{
		Port:           port,
		Environment:    env,
		RiotAPIKey:     riotKey,
		GovernorMode:   govMode,
		DatabaseURL:    dbURL,
		RedisURL:       redisURL,
		OpenSearchURL:  os.Getenv("OPENSEARCH_URL"),
		JWTSecret:      jwtSecret,
		CurrentPatch:   "16.21",
		DDragonVersion: "16.21.1",
	}
}

func (c *Config) PortInt() int {
	p, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return p
}
