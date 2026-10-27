# D — Data Ingestion: Riot Data Ingestion Pipeline, Match Deduplication, Normalization & Schema Validation
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-D`  
**Related Work Orders:** `WO-RS-0301..WO-RS-0350`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **D — Data Ingestion** (`Riot Data Ingestion Pipeline, Match Deduplication, Normalization & Schema Validation`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### End-to-End Ingestion Pipeline
1. **Trigger:** User clicks `Update` on Summoner Profile or searches an uncached Riot ID.
2. **Deduplication Gate:** Query `match_ingestion_status` for returned Match IDs (`SELECT match_id FROM match_ingestion_status WHERE match_id = ANY($1) AND status = 'COMPLETED'`). Only unindexed Match IDs are queued.
3. **Remake & Invalid Match Filtering:**
   - Matches with `gameDuration < 300` seconds or `gameEndedInEarlySurrender = true` are stored in `matches` with `is_remake = true` so players can see them in their match history (styled grey as `Remake`), but they are **strictly excluded** from Win/Loss, KDA, and Champion Analytics aggregations.
4. **Atomic Transaction Write:** Each match is normalized and inserted across `matches`, `match_teams`, `match_participants`, `match_items`, `match_runes`, and `match_objectives` inside a single PostgreSQL transaction (`REPEATABLE READ`), followed by `match_ingestion_status.status = 'COMPLETED'`.

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
- **Primary Tables Involved:** `matches, match_participants, match_teams, match_items, match_runes, match_objectives, match_timeline_events, match_ingestion_status, ingestion_jobs, job_failures`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `POST /api/v1/summoners/{region}/{puuid}/refresh, GET /api/v1/ingestion/jobs/{jobId}`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Profile 'Update' Button State Machine (Idle -> Queued -> Fetching -> Completed / Cooldown)`
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
  3. All associated Work Orders (`WO-RS-0301..WO-RS-0350`) produce verifiable execution logs and Git commit SHAs.
