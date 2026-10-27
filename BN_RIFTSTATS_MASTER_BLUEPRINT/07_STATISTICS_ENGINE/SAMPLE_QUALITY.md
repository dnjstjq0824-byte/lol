# Sample Size Thresholds, Bias Disclosure & Data Quality Guardrails
**Document ID:** `BN-RIFTSTATS-STAT-005`  

---

## 1. Mandatory Sample Quality Classification (`sample_quality`)
Every statistical row in BN RiftStats carries an explicit `sample_quality` enum based on sample size $N$:
1. **`HIGH_CONFIDENCE` ($N \ge 500$):** Rendered normally; eligible for `OP` / `Tier 1` ranking.
2. **`MODERATE` ($100 \le N < 500$):** Rendered normally; Wilson lower-bound shrinkage active.
3. **`LOW_SAMPLE` ($10 \le N < 100$):** Rendered with an amber `Low Sample (N={N})` warning pill; excluded from headline "Best Win Rate" callouts.
4. **`INSUFFICIENT_DATA` ($N < 10$):** Win rate is suppressed from default tables (hidden behind an explicit `"Show low-confidence samples (<10 games)"` toggle) so a `100% WR (1 game)` never misleads users.

## 2. Public Coverage & Bias Transparency
Because Riot's `MATCH-V5` API does not push a global firehose of all matches played worldwide, BN RiftStats harvests matches via user profile lookups and automated ladder seed crawlers (`LEAGUE-V4` / `LEAGUE-EXP-V4`). Every statistics page displays the exact **Total Analyzed Matches (`total_matches_in_bucket`)** and **Last Updated Timestamp** in the header—never claiming to represent 100% of all unindexed games on Riot servers.
