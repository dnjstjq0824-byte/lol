# Patch Isolation, Version Parsing & Rolling Window Policy
**Document ID:** `BN-RIFTSTATS-STAT-006`  

---

## 1. Strict Patch Version Extraction
From `MATCH-V5` `info.gameVersion` (e.g., `"16.21.618.4912"`), the ingestion normalizer extracts the canonical `major.minor` string (`"16.21"`).
- Matches from `"16.20"` and `"16.21"` are **never** blended in patch-specific views.
- During the first **6 hours** after a new patch release when a specific region/tier slice has fewer than `5,000` indexed matches, the UI displays an explicit banner: *"Patch 16.21 just launched ({N} matches analyzed so far). Switch to Patch 16.20 for full historical sample size."*
