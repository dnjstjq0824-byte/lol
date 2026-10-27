# A — Architecture: Global Match Analytics System Architecture, Modular Monolith + Worker Topology & Scalability
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-A`  
**Related Work Orders:** `WO-RS-0101..WO-RS-0150`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **A — Architecture** (`Global Match Analytics System Architecture, Modular Monolith + Worker Topology & Scalability`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Core Architectural Decisions
1. **Modular Go Backend + Dedicated Ingestion/Analytics Workers:**
   - `cmd/api`: Stateless HTTP/REST API server (`Go 1.23`, `Chi`/`Echo` router, OpenAPI 3.1 validation middleware, OpenTelemetry tracing).
   - `cmd/ingestor`: Asynchronous worker pool consuming Redis Streams / Asynq queues (`ingest:match:p0`, `ingest:match:p1`, `ingest:timeline:p2`) governed by the centralized Riot API Rate Limit Governor.
   - `cmd/aggregator`: Scheduled & incremental statistical aggregation engine computing 7-dimensional champion, build, rune, item, matchup, and synergy cubes in PostgreSQL (with zero-downtime promotion and read-replica routing).
2. **Data Tier Separation:**
   - **Raw Payload Archive:** Compressed (`zstd`) raw `MATCH-V5` JSON payloads archived to S3-compatible object storage (`s3://bn-riftstats-raw-matches/{patch}/{region}/{matchId}.json.zst`) for replay/re-aggregation without re-calling Riot APIs.
   - **Normalized Relational Store:** Partitioned PostgreSQL 16 tables (`matches`, `match_participants`, `match_teams`, `match_items`, `match_runes`, `match_objectives`, `match_timeline_events`).
   - **Pre-Computed OLAP Cubes:** `champion_stat_aggregates`, `champion_build_aggregates`, `rune_aggregates`, `item_aggregates`, `matchup_aggregates`, `role_aggregates`.
   - **Low-Latency Cache & Search:** Redis 7 cluster (rate limit buckets, singleflight locks, API response cache) + OpenSearch (Riot ID prefix/trigram autocomplete & fuzzy champion lookup).

### 2.1 Sequence & Processing Pipeline
```text
[Client / Trigger] 
   │
   ▼
[Nuxt 3 SSR / API Gateway (Go)] ──► [Input Validation & Auth/Rate Guard]
   │
   ├──► [L1/L2 Cache Lookup (Redis)] ──(Hit)──► [Return Validated DTO]
   │
   ▼ (Miss / Stale)
[Singleflight Coordinator] ──► [PostgreSQL / OpenSearch / Riot API Governor]
   │
   ▼
[Normalize & Persist] ──► [Publish Telemetry & Refresh Cache] ──► [Client Response]
```

## 3. Database Schema Mapping (`PostgreSQL 16`)
- **Primary Tables Involved:** `users, riot_accounts, matches, match_participants, champion_stat_aggregates, api_request_logs`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `/api/v1/health, /api/v1/regions, /api/v1/summoners/{region}/{riotId}, /api/v1/champions/tier-list`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Global Layout, Region Switcher, SSR Hydration Boundary, Stale-While-Revalidate Data Banner`
- **UI State Machine:**
  - `IDLE` → `LOADING` (skeleton shimmer with fixed dimensions, `CLS = 0`) → `SUCCESS` (rendered data + freshness badge) | `STALE_FALLBACK` (cached data + non-intrusive warning banner when Riot API is degraded) | `EMPTY` (clear zero-data explanation) | `ERROR` (actionable retry button + error code).

## 6. Error Handling, Security & Riot Policy Guardrails
- **Error Recovery:** Graceful degradation on upstream Riot `429`/`502`/`503`/`504` errors; never overwrites valid cached PostgreSQL rows with error payloads.
- **Security Controls:** Server-side secret isolation, strict RBAC, input sanitization, and parameterized queries.
- **Riot Policy Compliance:** Enforces zero Arena Augment/Item win-rate display, zero custom MMR/ELO calculation, zero spectator de-anonymization, and mandatory Riot legal attribution boilerplate.

## 7. Performance Targets, Testing & Acceptance Criteria
- **Performance SLA:** Cached API reads `P95 < 80ms`; uncached on-demand Riot ingestion `P95 < 1,800ms`; UI interaction responsiveness `INP < 100ms` across all 11 viewports (`320px`~`1920px`).
- **Verification Suite:** Unit tests (`go test -race`), OpenAPI schema conformance tests, PostgreSQL migration tests, and Playwright E2E tests mapped in `13_TESTING/MASTER_TEST_MATRIX.csv`.
- **Acceptance Criteria:**
  1. All mapped endpoints and UI states pass automated schema and E2E assertions with zero mocked bypasses in integration suites.
  2. Zero Riot API key leakage or policy violations detected by static and dynamic security scans.
  3. All associated Work Orders (`WO-RS-0101..WO-RS-0150`) produce verifiable execution logs and Git commit SHAs.
