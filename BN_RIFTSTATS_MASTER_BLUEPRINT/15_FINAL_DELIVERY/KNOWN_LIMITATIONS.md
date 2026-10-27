# Known Limitations, Policy Restrictions & Pre-Production Gates
**Document ID:** `BN-RIFTSTATS-DEL-003`  

1. **Riot Production API Key Gate (`BLOCKED` for Live Commercial Traffic until Riot Approval):** Development/Personal keys are limited to `20 req/1s` and `100 req/2m`. Public launch must wait for Riot Developer Relations approval of the registered prototype (`500 req/10s, 30,000 req/10m`).
2. **Riot Sign On (`RSO`) Gate:** RSO OAuth credentials are only issued after Production Key approval.
3. **Data Dragon Patch-Day Lag:** Data Dragon tarballs can lag behind live game client deployment by 1–6 hours; handled automatically via version fallback and placeholder icons.
4. **Historical Unindexed Matches:** Riot `MATCH-V5` retains match history for a rolling window (~2 years); matches older than Riot's retention window that were not previously ingested cannot be backfilled.
