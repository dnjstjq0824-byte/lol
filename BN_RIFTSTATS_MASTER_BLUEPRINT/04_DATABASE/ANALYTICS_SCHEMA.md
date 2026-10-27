# BN RiftStats — Analytical Cube & Materialized Aggregation Design
**Document ID:** `BN-RIFTSTATS-DB-003`  

---

## 1. Separation of Raw, Normalized, and Aggregated Data
1. **Raw Layer (`S3` + `match_ingestion_status`):** Stores untouched Riot `MATCH-V5` JSON compressed with `zstd`.
2. **Normalized OLTP Layer (`matches`, `match_participants`, `match_items`, `match_runes`):** Strictly validated relational rows where each match is ingested exactly once (`match_id` PK).
3. **Aggregated OLAP Cube Layer (`champion_stat_aggregates`, `champion_build_aggregates`, `rune_aggregates`, `item_aggregates`, `matchup_aggregates`):** Updated incrementally in micro-batches every `5 minutes` (and full-reconciled every `60 minutes`) so user queries on `/champions/tier-list` or `/champions/{slug}/build` execute single-table index scans in `< 5ms` without scanning millions of `match_participants` rows at request time.
