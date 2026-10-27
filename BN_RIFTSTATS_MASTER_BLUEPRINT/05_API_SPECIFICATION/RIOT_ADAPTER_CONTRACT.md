# Go Riot API Adapter Internal Contract
**Document ID:** `BN-RIFTSTATS-API-ADAPTER`  

Defines the internal Go interface `riotclient.Client` wrapping all 10 Riot API families (`ACCOUNT-V1`, `SUMMONER-V4`, `MATCH-V5`, `LEAGUE-V4`, `LEAGUE-EXP-V4`, `CHAMPION-MASTERY-V4`, `CHAMPION-V3`, `SPECTATOR-V5`, `LOL-CHALLENGES-V1`, `LOL-STATUS-V4`) with mandatory `context.Context` deadlines, priority tags (`PriorityP0`, `PriorityP1`, `PriorityP2`), Redis Lua token-bucket reservations, and automatic `429 Retry-After` observance.
