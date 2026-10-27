# BN RiftStats — Entity-Relationship Diagram (ERD) & Domain Relationships
**Document ID:** `BN-RIFTSTATS-DB-001`  

---

## 1. Relational Domain Overview (43 Canonical Tables)
```text
[IDENTITY DOMAIN]
users (1) ──── (1) user_profiles
  ├── (1..N) auth_identities (EMAIL / GOOGLE / DISCORD / RIOT_RSO)
  ├── (1..N) sessions
  ├── (1..N) favorites ──────────────────────────┐
  └── (1..N) notifications                       │
                                                 ▼
[RIOT ACCOUNTS & RANKING DOMAIN]                 │
account_regions (1) ──── (N) riot_accounts ◄─────┘
                              ├── (1..N) riot_id_history
                              ├── (1..N) summoner_profiles
                              ├── (1..N) ranked_entries
                              ├── (1..N) ranking_snapshots
                              ├── (1..N) leaderboard_entries
                              └── (1..N) champion_mastery_snapshots

[MATCHES & STATIC GAME DATA DOMAIN]
patches (1) ──── (N) matches (1) ──── (2) match_teams
champions (1) ─┐      ├── (10..16) match_participants
queues (1) ────┘      ├── (N) match_items
                      ├── (10) match_runes
                      ├── (N) match_objectives
                      ├── (N) match_timeline_events
                      └── (1) match_ingestion_status

[ANALYTICS & OPERATIONS DOMAIN]
champion_stat_aggregates / role_aggregates / champion_build_aggregates
rune_aggregates / item_aggregates / matchup_aggregates / analytics_jobs
api_request_logs / rate_limit_buckets / ingestion_jobs / job_failures
source_freshness / reconciliation_results / ad_placements / audit_logs
admin_actions / feature_flags / system_settings
```
