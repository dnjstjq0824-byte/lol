# Summoner, Account Lookup & Ranked Stats API Specification
**Document ID:** `BN-RIFTSTATS-API-SUMMONER`  

Covers API Groups 1, 2, 3, 4, 16, and 19 (`Regions`, `Riot Account Lookup`, `Summoner Profile`, `Ranked Stats`, `Champion Mastery`, and `Spectator`). Enforces `120s` per-summoner on-demand refresh cooldown (`POST /api/v1/ingestion/summoners/{region}/{puuid}/refresh`) and returns cached PostgreSQL data with `stale: true` if Riot API returns `5xx` or `429`.
