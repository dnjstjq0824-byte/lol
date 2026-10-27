# BN RiftStats — Multi-Region Deployment & Zero-Downtime Release Topology
**Document ID:** `BN-RIFTSTATS-ARCH-004`  

---

## 1. Environment Promotion Pipeline
1. **Local Sandbox (`docker-compose.dev.yml`):** Full local stack with PostgreSQL 16, Redis 7, OpenSearch, Mock Riot Server (for deterministic offline CI tests), and Nuxt/Go hot-reload.
2. **Staging / Riot Prototype Review Environment (`staging.riftstats.bn.internal`):** Connected to Riot Personal/Development API Key (`20 req/1s, 100 req/2m` governor profile) for Riot Developer Relations audit.
3. **Production Environment (`production`):** Activated **only after Riot Production API Key approval** (`500 req/10s, 30,000 req/10m` per region), running rolling blue/green deployments with automated database migration pre-checks.
