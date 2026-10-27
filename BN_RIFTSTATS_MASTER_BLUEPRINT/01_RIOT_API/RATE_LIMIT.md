# Riot API Rate Limit Governor, Priority Queue & Resilience Design
**Document ID:** `BN-RIFTSTATS-API-005`  
**Verified Against:** `https://developer.riotgames.com/docs/portal#_rate-limiting`  

---

## 1. Riot Rate Limit Headers & Semantics
Every Riot API HTTP response includes rate-limit telemetry headers that the Go Gateway parses and synchronizes into Redis (`rate_limit_buckets`):
- `X-App-Rate-Limit`: e.g. `500:10,30000:600` (500 requests per 10s, 30,000 requests per 600s for the region).
- `X-App-Rate-Limit-Count`: e.g. `42:10,1920:600` (current usage in window).
- `X-Method-Rate-Limit`: e.g. `2000:10` (endpoint-specific limit per region).
- `X-Method-Rate-Limit-Count`: e.g. `15:10`.
- `Retry-After`: Integer seconds to wait when HTTP `429 Too Many Requests` is returned (`X-Rate-Limit-Type: application | method | service`).

## 2. 20-Component Rate Limit Architecture (Article 5 Compliance)
1. **Central API Request Manager (`internal/riotclient/governor.go`):** All outbound Riot calls pass through a single Go service layer. Direct `http.Get` calls to `*.api.riotgames.com` are forbidden by linter rule.
2. **Secret Management:** API key loaded from Vault/Secrets Manager with hot-reload support on rotation.
3. **Per-Region App Buckets:** Separate Redis Lua token buckets per `platform_id` (`kr`, `na1`, etc.) and `regional_id` (`asia`, `americas`, `europe`, `sea`).
4. **Application Limit Enforcement:** Dynamic synchronization with `X-App-Rate-Limit` headers (supports seamless transition from `20/1s, 100/120s` Dev key to `500/10s, 30000/600s` Production key without code changes).
5. **Method Limit Enforcement:** Independent token bucket per `(region, endpoint_group)` dynamically calibrated from `X-Method-Rate-Limit`.
6. **HTTP 429 Handler:** Distinguishes `X-Rate-Limit-Type`:
   - `application`: Pauses all requests for that region for `Retry-After` seconds.
   - `method`: Pauses only that specific endpoint in that region for `Retry-After` seconds.
   - `service`: Underlying Riot service overload; triggers exponential backoff with jitter (`base=2s, max=30s`).
7. **Priority Budget Allocation (3-Tier Traffic Shaping):**
   - **Priority P0 (Interactive User Search & Live Game - 65% capacity):** Reserved up to 100% of bucket capacity.
   - **Priority P1 (On-Demand Match History Refresh - 25% capacity):** Throttled when regional bucket usage exceeds `80%`.
   - **Priority P2 (Background Seed Crawler & Ladder Sync - 10% capacity):** Automatically paused whenever regional bucket usage exceeds `60%`.
8. **Singleflight Deduplication (`golang.org/x/sync/singleflight` + Redis Lock):** Concurrent requests for the same `(region, endpoint, resource_id)` share a single in-flight Riot API call.
9. **Request Budget per User Action:**
   - **Initial Summoner Search / Profile Load:** Maximum **4 Riot API calls** (`ACCOUNT-V1` + `SUMMONER-V4` + `LEAGUE-V4` + `MATCH-V5 IDs`), with lazy/queued ingestion of unindexed match details (`<= 10` new matches per page load).
10. **Circuit Breaker:** Opens after 20 consecutive `5xx` errors or network timeouts in a region for `30s`, serving cached PostgreSQL data with `stale: true` banner.
