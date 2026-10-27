# K — KPI & Match Metrics: Mathematical Definitions of KDA, KP%, CS/m, GPM, DPM, Vision Score, Laning Ratio & Rift Score
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-K`  
**Related Work Orders:** `WO-RS-0355..WO-RS-0395, WO-RS-0401..WO-RS-0430`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **K — KPI & Match Metrics** (`Mathematical Definitions of KDA, KP%, CS/m, GPM, DPM, Vision Score, Laning Ratio & Rift Score`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Exact Match Metric Formulas
- **KDA Ratio:** $\text{KDA} = \frac{K + A}{\max(1, D)}$ (When $D = 0$, displayed with a `Perfect` badge alongside the numeric $(K+A):1$ ratio).
- **Kill Participation (`KP%`):** $\text{KP} = \frac{K_p + A_p}{\max(1, \sum_{i \in \text{Team}} K_i)} \times 100\%$.
- **CS Per Minute (`CS/m`):** $\text{CS/m} = \frac{\text{totalMinionsKilled} + \text{neutralMinionsKilled}}{\text{gameDurationSeconds} / 60}$.
- **Gold Per Minute (`GPM`) & Damage Per Minute (`DPM`):** $\text{GPM} = \frac{\text{goldEarned}}{\text{gameDurationSeconds}/60}$, $\text{DPM} = \frac{\text{totalDamageDealtToChampions}}{\text{gameDurationSeconds}/60}$.
- **Laning Ratio (`Laning X:Y`):** Derived from 14-minute timeline gold/XP share between lane opponents: $\text{LaningShare}_A = \text{round}\left(50 + 50 \cdot \tanh\left(\frac{\text{GD}@14 + 0.8 \cdot \text{XPD}@14}{2500}\right)\right)$.
- **Rift Score (`0.0` ~ `10.0`, Rank `1st`~`10th`, `MVP` / `ACE`):** Role-normalized weighted linear combination of KP%, DPM share, Damage Taken/Mitigated share, Vision Score/min, Gold Efficiency, Objective Participation, and Death Penalty. Explicitly labeled as a **Single-Match Performance Index** (never an MMR replacement).

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
- **Primary Tables Involved:** `match_participants, match_teams, match_timeline_events`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `GET /api/v1/summoners/{region}/{puuid}/matches, GET /api/v1/matches/{region}/{matchId}`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Match Summary Cards, 10-Player Scoreboard, Radar/Bar Comparison Charts, Rift Score Badge (1st~10th)`
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
  3. All associated Work Orders (`WO-RS-0355..WO-RS-0395, WO-RS-0401..WO-RS-0430`) produce verifiable execution logs and Git commit SHAs.
