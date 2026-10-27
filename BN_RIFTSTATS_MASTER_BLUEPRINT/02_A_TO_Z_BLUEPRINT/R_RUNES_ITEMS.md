# R — Runes, Items & Builds: Rune Tree Aggregation, Starter/Boots/Core/Depth 4-6 Item Builds, Skill Order & Summoner Spells
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-R`  
**Related Work Orders:** `WO-RS-0501..WO-RS-0550`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **R — Runes, Items & Builds** (`Rune Tree Aggregation, Starter/Boots/Core/Depth 4-6 Item Builds, Skill Order & Summoner Spells`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Build & Item Path Aggregation Logic
- **Starter Items:** Extracted from `MATCH-V5 Timeline` frame `t <= 90s` item purchases (`ITEM_PURCHASED` minus `ITEM_UNDO`).
- **Boots:** Extracted from final completed Tier-2/Tier-3 boot IDs (including Season 2026 upgraded boots).
- **Core 3-Item Builds:** Extracted from the chronological purchase order of the first 3 completed Legendary/Core items in `match_timeline_events` (fallback to inventory slot inspection if timeline is absent).
- **Depth 4 / 5 / 6 Situational Items:** Conditional win rate and pick rate of the 4th, 5th, and 6th completed legendary items given completion of a 3-item core.
- **Arena Policy Enforcement:** Any match with `queue_id IN (1700, 1710)` has Augment and Item win-rate aggregation hard-disabled at the SQL/Worker layer.

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
- **Primary Tables Involved:** `champion_build_aggregates, rune_aggregates, item_aggregates, match_items, match_runes`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `GET /api/v1/champions/{championKey}/builds, GET /api/v1/champions/{championKey}/runes, GET /api/v1/champions/{championKey}/items, GET /api/v1/champions/{championKey}/skills`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Champion Build Page (PG-PUB-14), Runes Tab (PG-PUB-15), Items Tab (PG-PUB-16), Skills Tab (PG-PUB-17)`
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
  3. All associated Work Orders (`WO-RS-0501..WO-RS-0550`) produce verifiable execution logs and Git commit SHAs.
