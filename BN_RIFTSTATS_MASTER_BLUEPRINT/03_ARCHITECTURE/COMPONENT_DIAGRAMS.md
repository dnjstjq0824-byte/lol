# BN RiftStats — C4 Component & Topology Diagrams
**Document ID:** `BN-RIFTSTATS-ARCH-002`  

---

## 1. System Context & Container Topology
```text
                        ┌──────────────────────────────────────────────┐
                        │           Global Players & Analysts          │
                        │    (Desktop 1024..1920px / Mobile 320..430px)│
                        └──────────────────────┬───────────────────────┘
                                               │ HTTPS (TLS 1.3)
                                               ▼
                        ┌──────────────────────────────────────────────┐
                        │         Global CDN & WAF (Cloudflare)        │
                        │   - Static Assets (/ddragon/*, JS/CSS)       │
                        │   - Edge Caching & DDoS Protection           │
                        └──────────┬────────────────────────┬──────────┘
                                   │                        │
                 ┌─────────────────▼──────┐        ┌────────▼────────────────┐
                 │  User Web (Nuxt 3 SSR) │        │ Admin Console (Vue 3)   │
                 │  - Public Pages (1-23) │        │ - Ops Dashboard (1-18)  │
                 └─────────────────┬──────┘        └────────┬────────────────┘
                                   │                        │
                                   └───────────┬────────────┘
                                               │ Internal REST (OpenAPI 3.1)
                                               ▼
                        ┌──────────────────────────────────────────────┐
                        │           Go Backend API Cluster             │
                        │  - Auth/RBAC Middleware                      │
                        │  - Summoner / Match / Champion / Live APIs   │
                        │  - Singleflight Request Coalescer            │
                        └────┬─────────────────┬──────────────────┬────┘
                             │                 │                  │
        ┌────────────────────▼───┐   ┌─────────▼─────────┐   ┌────▼───────────────────┐
        │   Redis 7.2 Cluster    │   │   PostgreSQL 16   │   │    OpenSearch 2.x      │
        │ - Riot Token Buckets   │   │ - Primary + Read  │   │ - Riot ID Autocomplete │
        │ - Hot API Cache        │   │ - 43 Tables       │   │ - Champion Alias Index │
        │ - Asynq Job Streams    │   │ - Partitioned     │   └────────────────────────┘
        └───────────┬────────────┘   └─────────▲─────────┘
                    │                          │
        ┌───────────▼──────────────────────────┴─────────┐
        │     Go Ingestion & Aggregation Worker Pool     │
        │  - Riot API Rate Limit Governor (P0 / P1 / P2) │
        │  - Match & Timeline Normalizer                 │
        │  - 7D Champion & Build Cube Aggregator         │
        └───────────┬──────────────────────────┬─────────┘
                    │                          │
        ┌───────────▼────────────┐   ┌─────────▼─────────┐
        │   Riot Games Official  │   │  S3 Raw Archive   │
        │   APIs & Data Dragon   │   │ (zstd Match JSON) │
        └────────────────────────┘   └───────────────────┘
```
