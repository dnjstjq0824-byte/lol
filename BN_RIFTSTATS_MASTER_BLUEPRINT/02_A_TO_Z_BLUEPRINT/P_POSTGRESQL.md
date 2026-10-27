# P — PostgreSQL: Complete 43-Table Relational Schema, Partitioning Strategy, Indexing, Materialized Views & Backups
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-P`  
**Related Work Orders:** `WO-RS-0801..WO-RS-0825`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **P — PostgreSQL** (`Complete 43-Table Relational Schema, Partitioning Strategy, Indexing, Materialized Views & Backups`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Schema Architecture & Partitioning
- **43 Canonical Tables** organized across 8 logical domains (full SQL DDL provided in `04_DATABASE/migrations/0001_initial_schema.up.sql`).
- **Declarative Range Partitioning:** `matches`, `match_participants`, `match_timeline_events`, and `api_request_logs` are partitioned by `patch_version` or monthly time ranges (`pg_partman`) so old patch data can be detached or archived to S3 without table bloat or vacuum stalls.

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
- **Primary Tables Involved:** `All 43 Tables across Identity, Riot Accounts, Ranking, Matches, Game Data, Analytics, Riot API Operations, and Platform`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `Database Layer for All Backend Services`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Admin Database Health, Partition Status & Reconciliation Console`
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
  3. All associated Work Orders (`WO-RS-0801..WO-RS-0825`) produce verifiable execution logs and Git commit SHAs.
