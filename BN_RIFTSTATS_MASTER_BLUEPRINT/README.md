# BN RiftStats — Global League of Legends Match History & Champion Analytics Platform
**Master Blueprint & 1,000-Work-Order Engineering Package (`v1.0.0`)**

---

## 1. Package Directory Structure
- `START_HERE.md` — 20-step autonomous agent startup & execution protocol
- `README.md` — Master architecture overview & directory guide
- `MANIFEST.json` — Machine-readable inventory of all files, byte sizes, and SHA-256 digests
- `FILE_CHECKSUMS.sha256` — Standard `sha256sum -c` compatible integrity file
- `00_MARKET_RESEARCH/` — OP.GG live functional audit (`Patch 16.21`), page matrix, 7-competitor 20-dimension benchmark, feature parity matrix, and market positioning
- `01_RIOT_API/` — Complete audit of 10 Riot APIs, endpoint matrix, Dev/Personal/Production key policies, 16-region routing table, Redis token-bucket rate limiter, RSO OAuth, Data Dragon `16.21.1`, and Riot policy compliance
- `02_A_TO_Z_BLUEPRINT/` — 26 standalone master specifications (`A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`)
- `03_ARCHITECTURE/` — System overview, C4 diagrams, 12-stage data ingestion pipeline, deployment topology, and STRIDE threat model
- `04_DATABASE/` — Complete PostgreSQL 16 ERD, 50-table data dictionary (covering all 43+ required tables), OLAP cube schema, index design, and executable SQL migrations (`0001_initial_schema.up.sql` & `down.sql`)
- `05_API_SPECIFICATION/` — Validated OpenAPI 3.1 specification (`openapi.yaml`) covering all 26 API groups, Riot adapter contract, and 120 contract test cases
- `06_FRONTEND_UXUI/` — Wireframes, component trees, state machines, and `RESPONSIVE_MATRIX.csv` for all 23 Public views, 9 User views, 18 Admin views, and 11 viewports (`320px`~`1920px`)
- `07_STATISTICS_ENGINE/` — Exact mathematical definitions for Win/Pick/Ban rates, Wilson Score Lower Bound Champion Tier Engine, Build/Rune/Item aggregation, Lane Counter & Duo Synergy engine, Sample Quality guards, and Patch isolation
- `08_RIOT_DATA_PIPELINE/` — Riot ID lookup, MATCH-V5 detail/timeline ingestion, regional routing, multi-tier caching, and hourly reconciliation specs
- `09_USER_ADMIN/` — Argon2id/MFA/RSO authentication, local+cloud favorites, 18-screen Admin Console, and CLS-safe ad operations
- `10_SECURITY_COMPLIANCE/` — Riot policy enforcement, Vault API key isolation, GDPR/PIPA privacy, IP/cleanroom rules, and security test cases
- `11_WORK_ORDERS_0001_TO_1000/` — `MASTER_WORK_ORDER.md`, `WORK_ORDER_INDEX.csv`, `DEPENDENCY_GRAPH.md`, and 1,000 individual Work Orders (`individual/WO-RS-0001.md` ~ `WO-RS-1000.md`), each with all 25 required fields and a 25-point implementation checklist
- `12_AGENT_EXECUTION/` — 16-agent orchestration directives, file ownership matrix, phased execution sequence, cross-agent audit rules, and reporting template
- `13_TESTING/` — `MASTER_TEST_MATRIX.csv`, Acceptance Criteria, and detailed E2E test specifications for Tests A through J
- `14_INFRASTRUCTURE/` — CI/CD pipelines, Docker/K8s deployment, OpenTelemetry/Prometheus monitoring, pgBackRest PITR backups, and patch-day runbooks
- `15_FINAL_DELIVERY/` — Phase 0–11 roadmap, requirements traceability matrix (`REQUIREMENTS_TRACEABILITY.csv`), feature parity report, known limitations, 20-point delivery checklist, and final master report
