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
		// Admin console is RBAC-protected: anonymous callers get 401 (see TestAdminRBAC).
		{"GET", "/api/v1/admin/overview", http.StatusUnauthorized},
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

func loginAs(t *testing.T, ts http.Handler, email, password string) string {
	t.Helper()
	body := strings.NewReader(`{"email":"` + email + `","password":"` + password + `"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Login as %s failed: %d, body: %s", email, rec.Code, rec.Body.String())
	}
	var authResp struct {
		AccessToken string `json:"accessToken"`
	}
	_ = json.NewDecoder(rec.Body).Decode(&authResp)
	if authResp.AccessToken == "" {
		t.Fatalf("Expected access token for %s", email)
	}
	return authResp.AccessToken
}

func TestAdminRBAC(t *testing.T) {
	ts := setupTestServer()

	// 1. Anonymous -> 401
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/admin/overview", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for anonymous admin access, got %d", rec.Code)
	}

	// 2. Admin token -> 200 with real telemetry shape
	adminToken := loginAs(t, ts, "admin@riftstats.bn.io", "securepassword")
	reqAdmin := httptest.NewRequest("GET", "/api/v1/admin/overview", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	recAdmin := httptest.NewRecorder()
	ts.ServeHTTP(recAdmin, reqAdmin)
	if recAdmin.Code != http.StatusOK {
		t.Fatalf("Expected 200 for admin overview, got %d: %s", recAdmin.Code, recAdmin.Body.String())
	}
	var overview struct {
		RateLimitBuckets []struct {
			PlatformID string `json:"platformId"`
		} `json:"rateLimitBuckets"`
		IngestionQueueDepth int `json:"ingestionQueueDepth"`
	}
	if err := json.NewDecoder(recAdmin.Body).Decode(&overview); err != nil {
		t.Fatalf("Failed decoding admin overview: %v", err)
	}

	// 3. Regular user -> 403
	regBody := strings.NewReader(`{"email":"user1@example.com","password":"userpass123","displayName":"User1"}`)
	reqReg := httptest.NewRequest("POST", "/api/v1/auth/register", regBody)
	reqReg.Header.Set("Content-Type", "application/json")
	recReg := httptest.NewRecorder()
	ts.ServeHTTP(recReg, reqReg)
	if recReg.Code != http.StatusCreated {
		t.Fatalf("Expected 201 for registration, got %d: %s", recReg.Code, recReg.Body.String())
	}
	userToken := loginAs(t, ts, "user1@example.com", "userpass123")
	reqUser := httptest.NewRequest("GET", "/api/v1/admin/overview", nil)
	reqUser.Header.Set("Authorization", "Bearer "+userToken)
	recUser := httptest.NewRecorder()
	ts.ServeHTTP(recUser, reqUser)
	if recUser.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 for non-admin overview access, got %d", recUser.Code)
	}
}

func TestAuthPasswordVerification(t *testing.T) {
	ts := setupTestServer()

	// Register a user, then verify wrong passwords are rejected (no auto-login).
	regBody := strings.NewReader(`{"email":"verify@example.com","password":"correct-horse-9","displayName":"Verify"}`)
	reqReg := httptest.NewRequest("POST", "/api/v1/auth/register", regBody)
	reqReg.Header.Set("Content-Type", "application/json")
	recReg := httptest.NewRecorder()
	ts.ServeHTTP(recReg, reqReg)
	if recReg.Code != http.StatusCreated {
		t.Fatalf("Expected 201 for registration, got %d: %s", recReg.Code, recReg.Body.String())
	}

	wrongBody := strings.NewReader(`{"email":"verify@example.com","password":"wrong-password"}`)
	reqWrong := httptest.NewRequest("POST", "/api/v1/auth/login", wrongBody)
	reqWrong.Header.Set("Content-Type", "application/json")
	recWrong := httptest.NewRecorder()
	ts.ServeHTTP(recWrong, reqWrong)
	if recWrong.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for wrong password, got %d", recWrong.Code)
	}

	// Unknown email must not auto-register.
	ghostBody := strings.NewReader(`{"email":"ghost@example.com","password":"whatever123"}`)
	reqGhost := httptest.NewRequest("POST", "/api/v1/auth/login", ghostBody)
	reqGhost.Header.Set("Content-Type", "application/json")
	recGhost := httptest.NewRecorder()
	ts.ServeHTTP(recGhost, reqGhost)
	if recGhost.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for unknown email, got %d", recGhost.Code)
	}

	// Duplicate registration -> 409.
	dupBody := strings.NewReader(`{"email":"verify@example.com","password":"another-pass-1"}`)
	reqDup := httptest.NewRequest("POST", "/api/v1/auth/register", dupBody)
	reqDup.Header.Set("Content-Type", "application/json")
	recDup := httptest.NewRecorder()
	ts.ServeHTTP(recDup, reqDup)
	if recDup.Code != http.StatusConflict {
		t.Fatalf("Expected 409 for duplicate registration, got %d", recDup.Code)
	}

	// Short password -> 400.
	shortBody := strings.NewReader(`{"email":"short@example.com","password":"123"}`)
	reqShort := httptest.NewRequest("POST", "/api/v1/auth/register", shortBody)
	reqShort.Header.Set("Content-Type", "application/json")
	recShort := httptest.NewRecorder()
	ts.ServeHTTP(recShort, reqShort)
	if recShort.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 for short password, got %d", recShort.Code)
	}
}

func TestFavoriteValidationAndRefreshCooldown(t *testing.T) {
	ts := setupTestServer()
	token := loginAs(t, ts, "admin@riftstats.bn.io", "securepassword")

	// Invalid favorite payload -> 400.
	badFav := strings.NewReader(`{"targetType":"BOGUS","platformRegion":"KR"}`)
	reqBad := httptest.NewRequest("POST", "/api/v1/favorites", badFav)
	reqBad.Header.Set("Content-Type", "application/json")
	reqBad.Header.Set("Authorization", "Bearer "+token)
	recBad := httptest.NewRecorder()
	ts.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 for invalid favorite, got %d", recBad.Code)
	}

	// Valid favorite -> 201 and listed back.
	goodFav := strings.NewReader(`{"targetType":"CHAMPION","platformRegion":"KR","targetChampionId":103}`)
	reqGood := httptest.NewRequest("POST", "/api/v1/favorites", goodFav)
	reqGood.Header.Set("Content-Type", "application/json")
	reqGood.Header.Set("Authorization", "Bearer "+token)
	recGood := httptest.NewRecorder()
	ts.ServeHTTP(recGood, reqGood)
	if recGood.Code != http.StatusCreated {
		t.Fatalf("Expected 201 for valid favorite, got %d: %s", recGood.Code, recGood.Body.String())
	}

	// Refresh cooldown coalescing: two rapid calls return the same job.
	var first, second struct {
		JobID                string `json:"jobId"`
		CooldownRemainingSec int    `json:"cooldownRemainingSec"`
	}
	for i, target := range []*struct {
		JobID                string `json:"jobId"`
		CooldownRemainingSec int    `json:"cooldownRemainingSec"`
	}{&first, &second} {
		_ = i
		req := httptest.NewRequest("POST", "/api/v1/ingestion/summoners/KR/puuid-cooldown-test/refresh", nil)
		rec := httptest.NewRecorder()
		ts.ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("Expected 202 for refresh, got %d", rec.Code)
		}
		if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
			t.Fatalf("Failed decoding refresh response: %v", err)
		}
	}
	if first.JobID == "" || first.JobID != second.JobID {
		t.Fatalf("Expected coalesced refresh job IDs, got %q and %q", first.JobID, second.JobID)
	}
	if second.CooldownRemainingSec <= 0 || second.CooldownRemainingSec > 120 {
		t.Fatalf("Expected cooldown within (0,120], got %d", second.CooldownRemainingSec)
	}
}

func TestSecurityHeadersPresent(t *testing.T) {
	ts := setupTestServer()
	req := httptest.NewRequest("GET", "/api/v1/regions", nil)
	rec := httptest.NewRecorder()
	ts.ServeHTTP(rec, req)
	for _, h := range []string{"X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy"} {
		if rec.Header().Get(h) == "" {
			t.Fatalf("Expected security header %s to be set", h)
		}
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
