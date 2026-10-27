# C — Champion Analytics: Champion Win Rate, Pick Rate, Ban Rate, Composite Tier Engine & Meta Delta Analytics
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-C`  
**Related Work Orders:** `WO-RS-0401..WO-RS-0500`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **C — Champion Analytics** (`Champion Win Rate, Pick Rate, Ban Rate, Composite Tier Engine & Meta Delta Analytics`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### 7-Dimensional Aggregation Cube
Every row in `champion_stat_aggregates` is uniquely keyed by:
`(patch_version, platform_region, tier_bracket, queue_id, position, champion_id, time_window)`
- **Supported `tier_bracket` Values:** `ALL`, `CHALLENGER`, `GRANDMASTER`, `MASTER_PLUS`, `MASTER`, `DIAMOND_PLUS`, `DIAMOND`, `EMERALD_PLUS` (Default), `EMERALD`, `PLATINUM_PLUS`, `PLATINUM`, `GOLD_PLUS`, `GOLD`, `SILVER`, `BRONZE`, `IRON`.
- **Supported `position` Values:** `ALL`, `TOP`, `JUNGLE`, `MIDDLE`, `BOTTOM`, `UTILITY` (mapped from Riot `teamPosition` / `individualPosition` with fallback inference confidence score).
- **Composite Tier Assignment (`OP`, `1`, `2`, `3`, `4`, `5`):**
  $$\text{TierScore}_{c,r} = 0.55 \cdot Z(\text{WR}_{\text{Wilson}}(c,r)) + 0.30 \cdot Z(\ln(1 + \text{PR}(c,r))) + 0.15 \cdot Z(\text{BR}(c,r)) \cdot \mathbb{I}(\text{WR}(c,r) \ge 0.50)$$
  where $\text{WR}_{\text{Wilson}}(c,r)$ is the lower bound of the 95% Wilson score confidence interval, preventing low-sample champions (`N < 100`) from artificially inflating into `OP` or `Tier 1`.

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
- **Primary Tables Involved:** `champion_stat_aggregates, role_aggregates, patches, champions, queues`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `GET /api/v1/champions/tier-list, GET /api/v1/champions/{championKey}/stats, GET /api/v1/champions/{championKey}/trends`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Champion Tier Table (PG-PUB-12), Role Filter Tabs, Tier/Patch/Region Selectors, Sample Confidence Badges`
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
  3. All associated Work Orders (`WO-RS-0401..WO-RS-0500`) produce verifiable execution logs and Git commit SHAs.
