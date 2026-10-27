# Automated Data Integrity & Statistical Reconciliation Jobs
**Document ID:** `BN-RIFTSTATS-PIPE-006`  

Hourly background verification job checking: (1) zero duplicate `match_id` counts, (2) 10 participants per Summoner's Rift match, (3) $\sum_{p} \text{RoleShare}_{c,p} = 100\%$, and (4) zero Arena Augment/Item win-rate rows in `item_aggregates`. Results are persisted in `reconciliation_results`.
