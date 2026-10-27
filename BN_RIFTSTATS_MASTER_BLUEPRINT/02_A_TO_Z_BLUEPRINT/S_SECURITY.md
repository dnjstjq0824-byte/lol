# S — Security: Riot API Key Isolation, OWASP ASVS L2 Controls, RBAC, CSRF/XSS/SQLi Defense & Privacy Compliance
**Blueprint Document ID:** `BN-RIFTSTATS-AZ-S`  
**Related Work Orders:** `WO-RS-0055..WO-RS-0075, WO-RS-0155..WO-RS-0195, WO-RS-0935..WO-RS-0950`  
**Compliance Standard:** Riot Games General Policies & LoL Game Integrity Policy  

---

## 1. Product Requirements & Functional Scope
This specification defines the complete engineering contract for **S — Security** (`Riot API Key Isolation, OWASP ASVS L2 Controls, RBAC, CSRF/XSS/SQLi Defense & Privacy Compliance`) within **BN RiftStats**, ensuring 100% functional parity with OP.GG's core League of Legends workflows while maintaining an independent technical and visual architecture.

## 2. Technical Implementation & Core Processing Flow
### Zero-Trust Security Controls
1. **Zero Client-Side Riot Key Exposure:** Confirmed by automated build-time secret scanner (`gitleaks` + Playwright network interceptor verifying no `X-Riot-Token` or `RGAPI-` string ever appears in browser traffic or bundles).
2. **Parameterized SQL Only:** 100% of PostgreSQL queries use `pgx` positional parameters (`$1, $2`); dynamic sort columns are validated against a strict Go `map[string]string` allowlist.
3. **Admin RBAC & Audit Trail:** Every state-changing action in `/api/v1/admin/*` requires `ROLE_SUPER_ADMIN` or `ROLE_OPS_ADMIN` + active MFA session, writing an immutable before/after JSON diff to `admin_actions` and `audit_logs`.

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
- **Primary Tables Involved:** `users, sessions, audit_logs, admin_actions, api_request_logs`
- **Transactional & Indexing Rules:** All writes use explicit parameterized transactions (`pgx/v5`), idempotent `ON CONFLICT DO UPDATE` upserts, and composite B-Tree/GIN indexes defined in `04_DATABASE/migrations/0001_initial_schema.up.sql`.

## 4. API Contract Mapping (`OpenAPI 3.1`)
- **Primary Endpoints:** `Auth, User, and Admin APIs`
- **Transport & Validation:** JSON over HTTPS (`TLS 1.3`), validated against `05_API_SPECIFICATION/openapi.yaml`, returning `application/problem+json` (`RFC 9457`) on any `4xx`/`5xx` condition.

## 5. UX/UI Integration & State Transitions
- **Mapped Views & Components:** `Security Headers (CSP, HSTS), Admin Security Audit Log Viewer`
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
  3. All associated Work Orders (`WO-RS-0055..WO-RS-0075, WO-RS-0155..WO-RS-0195, WO-RS-0935..WO-RS-0950`) produce verifiable execution logs and Git commit SHAs.
