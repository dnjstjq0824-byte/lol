# MATCH-V5 Idempotent Discovery, Fetching & Normalization Pipeline
**Document ID:** `BN-RIFTSTATS-PIPE-002`  

Implements `GET /lol/match/v5/matches/by-puuid/{puuid}/ids` discovery, deduplication against `match_ingestion_status`, priority-queued fetching (`P0` interactive vs `P2` crawler), remake detection (`gameDuration < 300`), and atomic PostgreSQL persistence + `zstd` raw S3 archiving.
