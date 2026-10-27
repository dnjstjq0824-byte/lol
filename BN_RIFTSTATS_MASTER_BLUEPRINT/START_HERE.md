# START HERE — BN RiftStats Master Blueprint & Autonomous Agent Execution Entrypoint
**Directive Code:** `BN-RIFTSTATS-GEMINI4-ARGON-OPGG-FULL-PARITY-1000-001`  
**Project Codename:** `BN RiftStats`  
**Target Benchmark:** `https://op.gg/` (100% Core LoL Functional Parity, Independent Brand & Architecture)  
**Observed Live Riot Patch:** `16.21` (Data Dragon `16.21.1`)  

---

## 1. Mandatory 20-Step Startup Protocol for Receiving Agents (Article 36)
Every autonomous engineering agent or human developer team receiving `BN_RIFTSTATS_OPGG_MASTER_BLUEPRINT_1000_V1_0.zip` MUST execute the following sequence in strict order:

1. **ZIP Integrity Verification:** Verify `FILE_CHECKSUMS.sha256` across all files in `BN_RIFTSTATS_MASTER_BLUEPRINT/`.
2. **Manifest Inspection:** Parse `MANIFEST.json` and confirm all 16 directories (`00_MARKET_RESEARCH` through `15_FINAL_DELIVERY`) and 1,100+ files are present.
3. **A~Z Blueprint Review:** Read all 26 specifications in `02_A_TO_Z_BLUEPRINT/A_ARCHITECTURE.md` through `Z_ZERO_TRUST.md`.
4. **Riot Policy Compliance Review:** Read `01_RIOT_API/RIOT_POLICY_COMPLIANCE.md` and `10_SECURITY_COMPLIANCE/RIOT_POLICY.md`. Confirm hard blocks on Arena Augment/Item win rates, custom MMR/ELO calculators, and spectator de-anonymization.
5. **Riot API Key Tier Check:** Inspect environment secret `RIOT_API_KEY` and set governor mode (`DEV_PROTOTYPE_MODE: 20 req/1s, 100 req/2m` vs `PRODUCTION_MODE: 500 req/10s, 30,000 req/10m`).
6. **Repository State Audit:** Check target Git repository and branch status.
7. **Gap Analysis:** Compare existing repository code against `11_WORK_ORDERS_0001_TO_1000/WORK_ORDER_INDEX.csv`.
8. **Dependency Graph Traversal:** Load `11_WORK_ORDERS_0001_TO_1000/DEPENDENCY_GRAPH.md` and assign Work Orders to the 16 specialized agents defined in `12_AGENT_EXECUTION/AGENT_ROLE_MATRIX.md`.
9. **Code Implementation:** Implement Go 1.23 backend, workers, and Vue 3 / Nuxt 3 + TypeScript frontend modules per `11_WORK_ORDERS_0001_TO_1000/individual/WO-RS-0001.md` ~ `WO-RS-1000.md`.
10. **Database Migration:** Apply `04_DATABASE/migrations/0001_initial_schema.up.sql` (50 canonical tables covering all 43+ required entities, indexes, and constraints).
11. **Riot API Gateway Integration:** Verify `ACCOUNT-V1`, `SUMMONER-V4`, `MATCH-V5`, `LEAGUE-V4`, `LEAGUE-EXP-V4`, `CHAMPION-MASTERY-V4`, `CHAMPION-V3`, `SPECTATOR-V5`, `LOL-CHALLENGES-V1`, and `LOL-STATUS-V4` adapters + Data Dragon `16.21.1` sync.
12. **User Web & Admin Console Wiring:** Connect all 23 Public views, 9 User views, and 18 Admin views (`06_FRONTEND_UXUI/`) to the 26 OpenAPI 3.1 endpoint groups (`05_API_SPECIFICATION/openapi.yaml`).
13. **Statistical Engine Execution:** Run 7D Champion, Build, Rune, Item, Counter, and Synergy aggregation workers (`07_STATISTICS_ENGINE/`).
14. **Real Data Verification:** Validate DB records and API outputs against actual/recorded Riot API payloads (zero fake/hardcoded win rates).
15. **Browser E2E Verification:** Run Playwright across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
16. **Security Verification:** Execute `13_TESTING/SECURITY_E2E.md` (zero API key leakage, RBAC, MFA, SQLi, XSS).
17. **Performance Verification:** Benchmark cached API `P95 < 80ms`, `CLS < 0.05`, `LCP < 1.8s`.
18. **Independent Zero-Trust Audit:** `AG-16 (Independent Audit Agent)` audits all completion evidence per `13_TESTING/INDEPENDENT_AUDIT.md`.
19. **Git Commit & Approved Release:** Record Git Commit SHAs per Work Order.
20. **Final Evidence Report:** Submit honest status classification (`VERIFIED` vs `BLOCKED` for live Production API Key / RSO approval gates).
