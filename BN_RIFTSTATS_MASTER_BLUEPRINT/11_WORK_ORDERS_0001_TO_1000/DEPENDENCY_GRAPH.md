# BN RiftStats — 1,000 Work Orders Topological Dependency Graph
**Document ID:** `BN-RIFTSTATS-WO-DAG`  
**Total Nodes:** `1,000` (`WO-RS-0001` ~ `WO-RS-1000`)  
**DAG Acyclicity:** Verified (`0 cycles`)  

## 1. Block-Level Phase & Dependency Flow
```text
[Block 01: WO-RS-0001..0050 OP.GG & Market Audit] ──┐
[Block 02: WO-RS-0051..0100 Riot API & Policy]    ──┼──► [Block 03: WO-RS-0101..0150 Core Architecture & Stack]
                                                    │        │
      ┌─────────────────────────────────────────────┘        ▼
      ├──► [Block 04: WO-RS-0151..0200 Identity & Auth] ◄────┤
      ├──► [Block 05: WO-RS-0201..0250 Riot ID & Routing] ◄──┤
      │         │                                            │
      │         ▼                                            │
      ├──► [Block 06: WO-RS-0251..0300 Summoner Profile]     │
      │         │                                            │
      │         ▼                                            │
      ├──► [Block 07: WO-RS-0301..0350 MATCH-V5 Ingestion] ◄─┴── [Block 17: WO-RS-0801..0850 DB/Redis/Queue/Infra]
      │         │
      │         ├──► [Block 08: WO-RS-0351..0400 Match History & Detail]
      │         ├──► [Block 09: WO-RS-0401..0450 Win/Pick/Ban Engine]
      │         │         │
      │         │         ├──► [Block 10: WO-RS-0451..0500 Tier Engine, Leaderboards & Patch]
      │         │         ├──► [Block 11: WO-RS-0501..0550 Runes, Items, Skills & Builds]
      │         │         └──► [Block 12: WO-RS-0551..0600 Counters & Synergies]
      │         └──► [Block 13: WO-RS-0601..0650 Game Modes & Spectator]
      │                   │
      ▼                   ▼
[Block 14: WO-RS-0651..0700 Home/Search/Profile UX] + [Block 15: WO-RS-0701..0750 Champion/Match/Stats UX]
      │
      ├──► [Block 16: WO-RS-0751..0800 Admin Console & Monetization]
      └──► [Block 18: WO-RS-0851..0900 Globalization, SEO & 11-Viewport Responsive]
                │
                ▼
          [Block 19: WO-RS-0901..0950 Unit, Integration, E2E & Security Tests]
                │
                ▼
          [Block 20: WO-RS-0951..1000 Independent Zero-Trust Audit & Commercial Gate]
```

## 2. Granular Block Entry & Exit Dependencies
- **Block 01 (`WO-RS-0001` ~ `WO-RS-0050` — OP.GG 전수 기능·시장 조사 (OP.GG Audit & Market Research)):** Assigned to `AG-02 (OP.GG Feature Audit Agent)` in `Phase 0`.
- **Block 02 (`WO-RS-0051` ~ `WO-RS-0100` — Riot API 정책·키·접근 권한 (Riot API Policy, Keys & Rate Governance)):** Assigned to `AG-03 (Riot API Integration Agent)` in `Phase 0`.
- **Block 03 (`WO-RS-0101` ~ `WO-RS-0150` — 시스템 아키텍처·개발 환경 (System Architecture & Core Dev Environment)):** Assigned to `AG-01 (Lead Architect) & AG-07 (Frontend Agent)` in `Phase 1`.
- **Block 04 (`WO-RS-0151` ~ `WO-RS-0200` — 회원가입·로그인·사용자 계정 (User Auth, Favorites & RSO Identity)):** Assigned to `AG-13 (Security Agent) & AG-14 (DevOps Agent)` in `Phase 2`.
- **Block 05 (`WO-RS-0201` ~ `WO-RS-0250` — Riot ID·PUUID·리전 검색 (Riot ID Resolution, PUUID & Global Region Routing)):** Assigned to `AG-03 (Riot API Integration Agent) & AG-12 (Search & SEO Agent)` in `Phase 2`.
- **Block 06 (`WO-RS-0251` ~ `WO-RS-0300` — 소환사 프로필·랭크·숙련도 (Summoner Profile, Ranked Stats & Mastery Engine)):** Assigned to `AG-09 (Summoner Profile Agent)` in `Phase 2`.
- **Block 07 (`WO-RS-0301` ~ `WO-RS-0350` — MATCH-V5 경기 수집 (MATCH-V5 Ingestion, Normalization & Deduplication)):** Assigned to `AG-04 (Data Ingestion Agent)` in `Phase 3`.
- **Block 08 (`WO-RS-0351` ~ `WO-RS-0400` — 전적 목록·경기 상세 (Match History Stream & 10-Player Match Detail Analysis)):** Assigned to `AG-11 (Match History Agent)` in `Phase 3`.
- **Block 09 (`WO-RS-0401` ~ `WO-RS-0450` — 챔피언 승률·픽률·밴률 (Champion Win Rate, Pick Rate & Ban Rate Engine)):** Assigned to `AG-05 (Game Statistics Agent)` in `Phase 5`.
- **Block 10 (`WO-RS-0451` ~ `WO-RS-0500` — 챔피언 티어·랭킹·패치 (Champion Tier Engine, Global Leaderboards & Patch Meta)):** Assigned to `AG-10 (Champion Analytics Agent)` in `Phase 5 & Phase 6`.
- **Block 11 (`WO-RS-0501` ~ `WO-RS-0550` — 룬·아이템·스킬·빌드 분석 (Runes, Items, Skill Order & Build Analytics)):** Assigned to `AG-10 (Champion Analytics Agent) & AG-05 (Game Statistics Agent)` in `Phase 5`.
- **Block 12 (`WO-RS-0551` ~ `WO-RS-0600` — 챔피언 카운터·시너지 (Champion Counters, Lane Matchups & Duo Synergies)):** Assigned to `AG-10 (Champion Analytics Agent) & AG-05 (Game Statistics Agent)` in `Phase 5`.
- **Block 13 (`WO-RS-0601` ~ `WO-RS-0650` — 게임 모드·Spectator 정책 대응 (Game Modes & SPECTATOR-V5 Policy Compliance)):** Assigned to `AG-09 (Summoner Profile Agent) & AG-03 (Riot API Integration Agent)` in `Phase 5`.
- **Block 14 (`WO-RS-0651` ~ `WO-RS-0700` — 사용자 홈·소환사 검색 UX/UI (Global Home, Search Flow & Summoner UX/UI)):** Assigned to `AG-08 (UX/UI Agent) & AG-07 (Frontend Agent)` in `Phase 7`.
- **Block 15 (`WO-RS-0701` ~ `WO-RS-0750` — 챔피언·전적·통계 UX/UI (Champion, Match Detail & Analytics UX/UI)):** Assigned to `AG-07 (Frontend Agent) & AG-08 (UX/UI Agent)` in `Phase 7`.
- **Block 16 (`WO-RS-0751` ~ `WO-RS-0800` — 관리자·광고·회원 운영 (Admin Console, Ad Operations & Monetization)):** Assigned to `AG-14 (DevOps Agent) & AG-06 (Database Agent)` in `Phase 9`.
- **Block 17 (`WO-RS-0801` ~ `WO-RS-0850` — DB·Redis·큐·검색·인프라 (PostgreSQL 43-Table DDL, Redis, Queue, OpenSearch & Infra)):** Assigned to `AG-06 (Database Agent) & AG-12 (Search & SEO Agent) & AG-14 (DevOps Agent)` in `Phase 1 & Phase 9`.
- **Block 18 (`WO-RS-0851` ~ `WO-RS-0900` — 글로벌화·SEO·성능·모바일 (Globalization, SEO, Core Web Vitals & 11-Viewport Mobile UX)):** Assigned to `AG-12 (Search & SEO Agent) & AG-08 (UX/UI Agent)` in `Phase 8`.
- **Block 19 (`WO-RS-0901` ~ `WO-RS-0950` — 단위·통합·E2E·보안 테스트 (Unit, Integration, E2E, Chaos & Security Verification)):** Assigned to `AG-15 (QA/E2E Agent) & AG-13 (Security Agent)` in `Phase 10`.
- **Block 20 (`WO-RS-0951` ~ `WO-RS-1000` — 독립 감사·상용화·배포 검증 (Independent Zero-Trust Audit, Production Approval & ZIP Packaging)):** Assigned to `AG-16 (Independent Audit Agent)` in `Phase 10 & Phase 11`.
