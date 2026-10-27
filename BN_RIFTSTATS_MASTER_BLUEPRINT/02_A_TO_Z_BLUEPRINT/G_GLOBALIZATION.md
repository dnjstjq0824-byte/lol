# G — Globalization: 16 Riot Regions, 23+ Locales, CJK/RTL Typography, Timezone Formatting & Regional SEO
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-G`  
**Related Work Orders:** `WO-RS-0851..WO-RS-0875`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **G — Globalization** (`16 Riot Regions, 23+ Locales, CJK/RTL Typography, Timezone Formatting & Regional SEO`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Multi-Region & Multi-Locale Engineering
- **Data Dragon Localization Mapping:** Automatically maps UI locale (`ko` → `ko_KR`, `en` → `en_US`, `ja` → `ja_JP`, `zh-CN` → `zh_CN`, `zh-TW` → `zh_TW`, `de` → `de_DE`, `es` → `es_ES`, `fr` → `fr_FR`, `pt` → `pt_BR`, `ru` → `ru_RU`, `tr` → `tr_TR`, `vi` → `vi_VN`, `th` → `th_TH`, `ar` → `ar_AE`) so champion names, item names, and rune descriptions render in the user's native language.
- **Cross-Language Search Alias Support:** Searching `"리신"` or `"Lee Sin"` or `"盲僧"` resolves to `champion_id = 64` (`LeeSin`) via `champion_localizations` trigram index.

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
- **Primary Tables Involved:** `account_regions, champion_localizations, patches, system_settings`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `GET /api/v1/regions, GET /api/v1/static/translations/{locale}`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Region Selector (PG-PUB-02), Locale Switcher (ko-KR, en-US, ja-JP, zh-CN, eu-DE, fr-FR, es-ES, pt-BR, tr-TR, vi-VN, ar-AE)`
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
  3. All associated Work Orders (`WO-RS-0851..WO-RS-0875`) produce verifiable execution logs and Git commit SHAs.
