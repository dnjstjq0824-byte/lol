# Match History, Match Detail & Match Timeline API Specification
**Document ID:** `BN-RIFTSTATS-API-MATCH`  

Covers API Groups 5, 6, and 7 (`Match History`, `Match Detail`, and `Match Timeline`). Uses keyset cursor pagination (`game_end_timestamp_ms + '_' + match_id`), filters out `is_private_custom = true` matches unless requested by the RSO-authenticated owner, and gracefully degrades if `MATCH-V5 Timeline` is missing.
