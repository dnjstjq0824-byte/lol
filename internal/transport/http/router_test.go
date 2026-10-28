package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bn-riftstats/internal/config"
	"bn-riftstats/internal/repository"
	"bn-riftstats/internal/riotclient"
	transport "bn-riftstats/internal/transport/http"
	"bn-riftstats/internal/transport/http/handlers"
)

func setupTestServer() http.Handler {
	cfg := &config.Config{
		Port:           "8080",
		Environment:    "test",
		GovernorMode:   "DEV_PROTOTYPE_MODE",
		JWTSecret:      "test-secret-key-1234567890",
		CurrentPatch:   "16.21",
		DDragonVersion: "16.21.1",
	}
	client := riotclient.NewClient(cfg)
	repo := repository.NewRepository(client.DDragon())
	h := handlers.NewHandler(cfg, client, repo)
	return transport.NewRouter(cfg, h)
}

func TestAll26EndpointGroups(t *testing.T) {
	ts := setupTestServer()

	endpoints := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/api/v1/regions", http.StatusOK},
		{"GET", "/api/v1/accounts/by-riot-id/KR/Hide%20on%20bush/KR1", http.StatusOK},
		{"GET", "/api/v1/search/autocomplete?q=faker", http.StatusOK},
		{"GET", "/api/v1/summoners/KR/by-puuid/puuid-test-123", http.StatusOK},
		{"GET", "/api/v1/summoners/KR/by-puuid/puuid-test-123/ranked", http.StatusOK},
		{"GET", "/api/v1/summoners/KR/by-puuid/puuid-test-123/matches", http.StatusOK},
		{"GET", "/api/v1/summoners/KR/by-puuid/puuid-test-123/mastery", http.StatusOK},
		{"GET", "/api/v1/matches/KR/KR_7123450000", http.StatusOK},
		{"GET", "/api/v1/matches/KR/KR_7123450000/timeline", http.StatusOK},
		{"GET", "/api/v1/champions", http.StatusOK},
		{"GET", "/api/v1/champions/tier-list", http.StatusOK},
		{"GET", "/api/v1/champions/Ahri/builds", http.StatusOK},
		{"GET", "/api/v1/champions/Ahri/runes", http.StatusOK},
		{"GET", "/api/v1/champions/Ahri/items", http.StatusOK},
		{"GET", "/api/v1/champions/Ahri/counters", http.StatusOK},
		{"GET", "/api/v1/champions/Ahri/synergies", http.StatusOK},
		{"GET", "/api/v1/leaderboards/KR?queue=420", http.StatusOK},
		{"GET", "/api/v1/modes/arena/stats", http.StatusOK},
		{"GET", "/api/v1/patches", http.StatusOK},
		{"GET", "/api/v1/spectator/KR/by-puuid/puuid-test-123", http.StatusOK},
		{"GET", "/api/v1/admin/overview", http.StatusOK},
		{"POST", "/api/v1/ingestion/summoners/KR/puuid-test-123/refresh", http.StatusAccepted},
		{"GET", "/api/v1/analytics/summary", http.StatusOK},
		{"GET", "/api/v1/monitoring/health", http.StatusOK},
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(ep.method, ep.path, nil)
		rec := httptest.NewRecorder()
		ts.ServeHTTP(rec, req)

		if rec.Code != ep.expectedStatus {
			t.Errorf("[%s %s] Expected status %d, got %d. Body: %s", ep.method, ep.path, ep.expectedStatus, rec.Code, rec.Body.String())
		}

		// Security Check: Verify zero API key leakage in responses
		body := rec.Body.String()
		if strings.Contains(body, "RGAPI-") {
			t.Fatalf("[%s %s] CRITICAL SECURITY VIOLATION: Riot API Key leaked in response body!", ep.method, ep.path)
		}
	}
}

func TestProblemDetailsErrorResponse(t *testing.T) {
	ts := setupTestServer()

	// Invalid region should produce RFC 9457 Problem Details
	req := httptest.NewRequest("GET", "/api/v1/accounts/by-riot-id/INVALID_REGION/player/tag", nil)
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 for invalid region, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatalf("Expected application/problem+json Content-Type, got %s", rec.Header().Get("Content-Type"))
	}

	var prob struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&prob); err != nil {
		t.Fatalf("Failed decoding problem details: %v", err)
	}
	if prob.Status != 404 || prob.Title != "Account Not Found" {
		t.Fatalf("Invalid problem payload: %+v", prob)
	}
}

func TestAuthFlowAndProtection(t *testing.T) {
	ts := setupTestServer()

	// 1. Login with valid credentials
	loginBody := strings.NewReader(`{"email":"admin@riftstats.bn.io","password":"securepassword"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", loginBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Login failed: %d, body: %s", rec.Code, rec.Body.String())
	}

	var authResp struct {
		AccessToken string `json:"accessToken"`
	}
	_ = json.NewDecoder(rec.Body).Decode(&authResp)
	if authResp.AccessToken == "" {
		t.Fatalf("Expected access token in login response")
	}

	// 2. Access /api/v1/users/me with token
	reqMe := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
	recMe := httptest.NewRecorder()
	ts.ServeHTTP(recMe, reqMe)

	if recMe.Code != http.StatusOK {
		t.Fatalf("Authorized /users/me failed: %d, body: %s", recMe.Code, recMe.Body.String())
	}

	// 3. Access without token should fail
	reqNoAuth := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	recNoAuth := httptest.NewRecorder()
	ts.ServeHTTP(recNoAuth, reqNoAuth)

	if recNoAuth.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for unauthenticated /users/me, got %d", recNoAuth.Code)
	}
}

func TestMultiSearchLobbyScouting(t *testing.T) {
	ts := setupTestServer()

	body := strings.NewReader(`{"region":"KR","inputs":["Hide on bush#KR1\nChovy#KR1\nShowMaker#KR1"]}`)
	req := httptest.NewRequest("POST", "/api/v1/search/multi-search", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Multi-search failed: %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Region  string `json:"region"`
		Members []struct {
			Account struct {
				GameName string `json:"gameName"`
			} `json:"account"`
			SoloQueue struct {
				Tier string `json:"tier"`
			} `json:"soloQueue"`
			TopChampions []struct {
				ChampionName string `json:"championName"`
			} `json:"topChampions"`
		} `json:"members"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed decoding multi-search response: %v", err)
	}

	if len(resp.Members) != 3 {
		t.Fatalf("Expected 3 members parsed, got %d", len(resp.Members))
	}
	if resp.Members[0].Account.GameName != "Hide on bush" {
		t.Fatalf("Expected member 0 to be Hide on bush, got %s", resp.Members[0].Account.GameName)
	}
}
