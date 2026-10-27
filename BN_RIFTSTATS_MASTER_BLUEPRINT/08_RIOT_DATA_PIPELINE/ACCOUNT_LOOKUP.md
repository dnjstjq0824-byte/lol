# Riot ID to PUUID Resolution & Historical Name Tracking Pipeline
**Document ID:** `BN-RIFTSTATS-PIPE-001`  

Handles Unicode `NFC` normalization, `#` splitting, default regional tag fallback, `ACCOUNT-V1` `/by-riot-id` lookup, `riot_accounts` upsert by immutable `puuid`, and automatic tracking of changed Riot IDs in `riot_id_history`.
