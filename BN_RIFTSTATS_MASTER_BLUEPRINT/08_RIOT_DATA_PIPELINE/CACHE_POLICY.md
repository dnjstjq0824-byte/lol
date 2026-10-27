# Multi-Tier Cache Hierarchy (CDN, Redis, PostgreSQL & Stale-While-Revalidate)
**Document ID:** `BN-RIFTSTATS-PIPE-005`  

Defines exact TTLs, cache keys, singleflight lock semantics, and stale-while-revalidate fallback rules so an upstream Riot API outage never causes a blank error screen for previously indexed summoners or champions.
