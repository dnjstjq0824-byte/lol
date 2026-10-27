# BN RiftStats — End-to-End Data Ingestion & Aggregation Pipeline
**Document ID:** `BN-RIFTSTATS-ARCH-003`  

---

## 1. 12-Stage Riot Data Ingestion Pipeline (Article 22)
```text
[1. Riot ID Search (gameName#tagLine)]
   │
   ▼
[2. Account Lookup (ACCOUNT-V1 /by-riot-id)] ──► Resolve canonical 78-char PUUID
   │
   ▼
[3. PUUID Resolution & History Tracking] ──► Upsert `riot_accounts` & `riot_id_history`
   │
   ▼
[4. Summoner Profile & Ranked Fetch] ──► Parallel `SUMMONER-V4` + `LEAGUE-V4` + `CHAMPION-MASTERY-V4`
   │
   ▼
[5. Match ID Discovery (MATCH-V5 /ids)] ──► Fetch latest 20 Match IDs via Regional Host
   │
   ▼
[6. Deduplication Filter] ──► Filter out Match IDs already `COMPLETED` in `match_ingestion_status`
   │
   ▼
[7. Match Detail Fetch (MATCH-V5 /matches/{id})] ──► Governed by P0/P1 Token Bucket
   │
   ▼
[8. Timeline Fetch (MATCH-V5 /timeline)] ──► Asynchronous P1/P2 worker fetch (non-blocking for basic card)
   │
   ▼
[9. Validation & Normalization] ──► Extract `major.minor` patch, remake flag, role confidence, 10 participants
   │
   ▼
[10. Atomic PostgreSQL Persistence + S3 Raw Archive] ──► Write to `matches`, `match_participants`, etc.
   │
   ▼
[11. Incremental Statistical Aggregation] ──► Update `champion_stat_aggregates`, `rune_aggregates`, `item_aggregates`, `matchup_aggregates` (excluding Remakes & Arena Augment/Item WRs)
   │
   ▼
[12. Cache Invalidation & User API Delivery] ──► Warm Redis keys & push updated payload to Nuxt UI
```
