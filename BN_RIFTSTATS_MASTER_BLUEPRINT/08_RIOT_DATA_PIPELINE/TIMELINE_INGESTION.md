# MATCH-V5 Timeline Event Extraction & Graceful Degradation
**Document ID:** `BN-RIFTSTATS-PIPE-003`  

Parses 60-second participant frames (`totalGold`, `xp`, `minionsKilled`, `jungleMinionsKilled`) to compute `@14m`/`@15m` lane differentials (`GD@15`, `XPD@15`, `CSD@15`), starter item bundles (`t <= 90s`), skill level-up order (`SKILL_LEVEL_UP`), and completed item purchase timestamps. If Timeline fetch fails or returns `404`, `timeline_status` is marked `NOT_AVAILABLE` without blocking the core Match Detail view.
