# 16-Agent Role, Responsibility & Exclusive File Ownership Matrix (Article 32)
**Document ID:** `BN-RIFTSTATS-AG-003`  

| Agent ID | Agent Role | Assigned Work Orders | Exclusive Directory / Codebase Ownership |
|---|---|---|---|
| `AG-01` | **Lead Architect** | `WO-RS-0101..0125` | `02_A_TO_Z_BLUEPRINT/`, `03_ARCHITECTURE/`, `05_API_SPECIFICATION/openapi.yaml` |
| `AG-02` | **OP.GG Feature Audit Agent** | `WO-RS-0001..0050` | `00_MARKET_RESEARCH/`, `15_FINAL_DELIVERY/FEATURE_PARITY_REPORT.md` |
| `AG-03` | **Riot API Integration Agent** | `WO-RS-0051..0100`, `WO-RS-0201..0250` | `01_RIOT_API/`, `internal/riotclient/`, `internal/domain/account/` |
| `AG-04` | **Data Ingestion Agent** | `WO-RS-0301..0350` | `08_RIOT_DATA_PIPELINE/`, `cmd/ingestor/`, `internal/domain/ingestion/` |
| `AG-05` | **Game Statistics Agent** | `WO-RS-0401..0450` | `07_STATISTICS_ENGINE/`, `cmd/aggregator/`, `internal/domain/analytics/` |
| `AG-06` | **Database Agent** | `WO-RS-0801..0825` | `04_DATABASE/`, `migrations/`, `internal/repository/postgres/` |
| `AG-07` | **Frontend Agent** | `WO-RS-0126..0150`, `WO-RS-0701..0750` | `web/src/components/`, `web/src/composables/`, `web/src/stores/` |
| `AG-08` | **UX/UI Agent** | `WO-RS-0651..0700`, `WO-RS-0876..0900` | `06_FRONTEND_UXUI/`, `web/src/assets/css/`, `web/src/layouts/` |
| `AG-09` | **Summoner Profile Agent** | `WO-RS-0251..0300`, `WO-RS-0601..0650` | `internal/domain/summoner/`, `internal/domain/spectator/`, `web/src/pages/summoners/` |
| `AG-10` | **Champion Analytics Agent** | `WO-RS-0451..0600` | `internal/domain/champion/`, `web/src/pages/champions/`, `web/src/pages/leaderboards/` |
| `AG-11` | **Match History Agent** | `WO-RS-0351..0400` | `internal/domain/match/`, `web/src/pages/matches/` |
| `AG-12` | **Search & SEO Agent** | `WO-RS-0830..0875` | `internal/repository/opensearch/`, `web/src/server/routes/sitemap.xml.ts` |
| `AG-13` | **Security Agent** | `WO-RS-0151..0200`, `WO-RS-0935..0950` | `10_SECURITY_COMPLIANCE/`, `internal/transport/http/middleware/auth.go` |
| `AG-14` | **DevOps Agent** | `WO-RS-0751..0800`, `WO-RS-0826..0829` | `14_INFRASTRUCTURE/`, `09_USER_ADMIN/`, `admin-web/`, `docker/`, `.github/workflows/` |
| `AG-15` | **QA/E2E Agent** | `WO-RS-0901..0934` | `13_TESTING/`, `e2e/playwright/`, `internal/testutil/` |
| `AG-16` | **Independent Audit Agent** | `WO-RS-0951..1000` | `12_AGENT_EXECUTION/`, `15_FINAL_DELIVERY/` |
