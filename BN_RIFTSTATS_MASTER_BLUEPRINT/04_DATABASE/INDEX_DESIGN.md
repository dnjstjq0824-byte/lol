# BN RiftStats — PostgreSQL Indexing, Partitioning & Backup Strategy
**Document ID:** `BN-RIFTSTATS-DB-004`  

---

## 1. High-Selectivity Index Design
- **Riot ID Exact & Prefix Search:** `idx_riot_accounts_name_tag` (`B-Tree` on `(game_name_normalized, tag_line_normalized)`) + `idx_riot_accounts_gamename_trgm` (`GIN` `pg_trgm` index) for `<15ms` fuzzy matching.
- **Summoner Match History Stream:** `idx_match_participants_puuid_match` (`(puuid, match_id DESC)`) enables instant keyset pagination for any player's match history.
- **Champion Tier Table Scan:** `idx_champ_stats_lookup` (`(patch_version, platform_region, tier_bracket, queue_id, position, tier_score DESC)`) serves the entire Champion Tier List in a single index-only scan.

## 2. Backup, Point-in-Time Recovery (PITR) & Disaster Recovery
- Continuous WAL archiving via `pgBackRest` to S3 with **RPO < 60 seconds** and **RTO < 15 minutes**, verified weekly via automated restore drills into an isolated verification container.
