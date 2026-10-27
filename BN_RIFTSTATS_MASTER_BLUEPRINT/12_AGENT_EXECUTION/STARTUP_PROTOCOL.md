# Receiving Agent 20-Step Startup & Execution Protocol (Article 36)
**Document ID:** `BN-RIFTSTATS-AG-002`  

1. Verify ZIP SHA-256 checksum (`FILE_CHECKSUMS.sha256`).
2. Validate all 1,105+ files against `MANIFEST.json`.
3. Read `02_A_TO_Z_BLUEPRINT/A_ARCHITECTURE.md` ~ `Z_ZERO_TRUST.md`.
4. Review `01_RIOT_API/RIOT_POLICY_COMPLIANCE.md`.
5. Check `RIOT_API_KEY` environment tier (`DEV_KEY` vs `PROD_KEY`).
6. Inspect existing Git repository state.
7. Run Gap Analysis against `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`.
8. Traverse `DEPENDENCY_GRAPH.md` topological order.
9. Implement Go / Nuxt / SQL code per Work Order spec.
10. Apply and verify PostgreSQL migrations (`0001_initial_schema.up.sql`).
11. Integrate Riot API Gateway with token-bucket governor.
12. Connect Nuxt 3 SSR/CSR views to OpenAPI 3.1 endpoints.
13. Run statistical aggregation workers and sample quality checks.
14. Validate live/recorded Riot data accuracy.
15. Execute Playwright E2E across all 11 viewports (`320px`~`1920px`).
16. Execute security test suite (`TEST-SEC-01..20`).
17. Execute load & latency benchmarks (`P95 < 80ms` cached).
18. Submit evidence package to `Independent Audit Agent`.
19. Record Git commit SHA and tag verified release.
20. Generate final audit report with honest `VERIFIED` / `BLOCKED` status classification.
