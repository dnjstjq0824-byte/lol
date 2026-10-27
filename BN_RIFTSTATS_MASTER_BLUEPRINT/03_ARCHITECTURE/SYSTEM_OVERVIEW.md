# BN RiftStats — Global System Architecture Overview
**Document ID:** `BN-RIFTSTATS-ARCH-001`  

---

## 1. Architectural Vision & Design Principles
**BN RiftStats** is architected as a high-throughput, low-latency global League of Legends analytics platform capable of serving millions of daily summoner lookups and multi-dimensional champion aggregations across all 16 Riot regions.

### 1.1 Technology Stack Summary (Article 21 Mandatory Stack)
| Layer | Technology | Role & Justification |
|---|---|---|
| **User Web Frontend** | `Vue 3.5` + `Nuxt 3.13` + `TypeScript 5.6` + `Tailwind CSS` + `Pinia` | Hybrid SSR/ISR for SEO-critical Summoner and Champion pages + fast CSR transitions |
| **Admin Console** | `Vue 3.5` + `TypeScript 5.6` + `Vite` | Dedicated RBAC-protected operational dashboard for Riot API telemetry, ingestion queues, and data quality |
| **Backend API & Workers** | `Go 1.23` (`chi/v5`, `pgx/v5`, `asynq`, `singleflight`) | High-concurrency, low-memory-footprint API gateway, Riot rate-limit governor, and statistical aggregation workers |
| **Primary Database** | `PostgreSQL 16` (Partitioned tables + Read Replicas) | ACID storage for accounts, normalized matches, rankings, and pre-computed statistical cubes |
| **Cache & Rate Coordination** | `Redis 7.2 Cluster` | Distributed Lua token buckets for Riot API rate limits, singleflight locks, session store, and response caching |
| **Search & Autocomplete** | `OpenSearch 2.x` (with PostgreSQL `pg_trgm` fallback) | `<50ms` prefix and fuzzy autocomplete for Riot IDs (`gameName#tagLine`) and localized champion names |
| **Raw Payload & Asset Storage** | `S3-Compatible Object Storage` + `Global CDN` | Permanent compressed archive (`zstd`) of raw `MATCH-V5` JSON payloads + mirrored Data Dragon static assets |
| **Observability** | `OpenTelemetry` + `Prometheus` + `Grafana` + `Loki` | End-to-end distributed tracing, Riot rate-limit utilization gauges, and ingestion latency histograms |
| **Testing & QA** | `Playwright` (Browser E2E) + `Go Test` (Unit/Integration/Contract) | Automated verification across all 11 viewports (`320px`~`1920px`) and 26 API groups |
