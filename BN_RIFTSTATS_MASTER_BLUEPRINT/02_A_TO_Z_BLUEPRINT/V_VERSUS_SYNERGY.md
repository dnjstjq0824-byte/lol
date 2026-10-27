# V — Versus & Synergy: Lane Counter Engine, Position Confidence Scoring, Matchup @15m Metrics & Team Duo Synergy Matrix
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-V`  
**Related Work Orders:** `WO-RS-0551..WO-RS-0600`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **V — Versus & Synergy** (`Lane Counter Engine, Position Confidence Scoring, Matchup @15m Metrics & Team Duo Synergy Matrix`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Counter & Synergy Statistical Engine
- **Direct Lane Matchup Definition:** Champion $A$ in role $R$ vs. Champion $B$ in the same role $R$ on the opposing team (`team_id_A != team_id_B` and `position_A == position_B`).
- **Position Confidence Guard:** Uses Riot's `teamPosition` when `position_confidence >= 0.85`. If a lane swap occurred or `teamPosition` is ambiguous, the match is flagged with `low_position_confidence = true` and excluded from strict lane-phase `@15m` counter diffs.
- **Duo Synergy Definition:** Champion $A$ in role $R_1$ paired with Champion $B$ in role $R_2$ ($R_1 \neq R_2$) on the **same team** (`team_id_A == team_id_B`). Computes combined Win Rate and **Synergy Delta**:
  $$\Delta_{\text{Synergy}}(A, B) = \text{WR}(A \land B) - \mathbb{E}[\text{WR}(A, B)]$$

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
- **Primary Tables Involved:** `matchup_aggregates, role_aggregates, match_participants`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `GET /api/v1/champions/{championKey}/counters, GET /api/v1/champions/{championKey}/synergies`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Champion Counters Page (PG-PUB-18), Champion Synergies Page (PG-PUB-19)`
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
  3. All associated Work Orders (`WO-RS-0551..WO-RS-0600`) produce verifiable execution logs and Git commit SHAs.
