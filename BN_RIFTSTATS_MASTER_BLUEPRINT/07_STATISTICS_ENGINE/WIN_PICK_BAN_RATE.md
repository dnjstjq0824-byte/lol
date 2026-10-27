# Mathematical Specification: Win Rate, Pick Rate, Ban Rate & Role Distribution
**Document ID:** `BN-RIFTSTATS-STAT-001`  

---

## 1. Aggregation Universe & Filtering Rules
Let $\mathcal{M}(v, r, b, q, w)$ denote the set of deduplicated, completed matches in PostgreSQL `matches` matching:
- Patch Version $v$ (`major.minor`, e.g., `16.21`) — **Strict Patch Isolation:** Matches from `16.20` are NEVER mixed into `16.21` default aggregations.
- Platform Region $r \in \{\text{KR}, \text{NA1}, \dots, \text{GLOBAL}\}$.
- Rank Tier Bracket $b$ (determined by the match's average visible participant rank in `ranked_entries` at ingestion time).
- Queue ID $q$ (e.g., `420` for Ranked Solo/Duo).
- Time Window $w$.
- **Remake & Invalid Match Exclusion:** Any match with `is_remake = TRUE` (`game_duration_sec < 300` or `gameEndedInEarlySurrender = TRUE`) or `is_private_custom = TRUE` is strictly excluded from $\mathcal{M}$.

Let $M = |\mathcal{M}(v, r, b, q, w)|$ be the total number of valid matches in the bucket.

## 2. Exact Formulas (Article 23 Compliance)
1. **Champion Role Win Rate ($\text{WR}_{c,p}$):**
   $$\text{WR}_{c,p} = \frac{\text{Wins}(c, p)}{\text{Picks}(c, p)} \times 100\%$$
   - **Numerator:** Number of participant records in $\mathcal{M}$ where `champion_id = c`, `team_position = p`, and `win = TRUE`.
   - **Denominator:** Total number of participant records in $\mathcal{M}$ where `champion_id = c` and `team_position = p`.

2. **Champion Role Pick Rate ($\text{PR}_{c,p}$):**
   $$\text{PR}_{c,p} = \frac{\text{Picks}(c, p)}{M} \times 100\%$$
   - **Numerator:** Number of participant records in $\mathcal{M}$ where `champion_id = c` and `team_position = p` (since a champion can be picked at most once per match in draft queues).
   - **Denominator:** Total valid matches $M$ in the bucket (so $\sum_{c} \text{PR}_{c,p} = 200\%$ across both teams in a single role $p$).

3. **Champion Ban Rate ($\text{BR}_{c}$):**
   $$\text{BR}_{c} = \frac{\text{MatchesWithBan}(c)}{M} \times 100\%$$
   - **Numerator:** Number of distinct matches $m \in \mathcal{M}$ where `c = ANY(team_100.bans_champion_ids || team_200.bans_champion_ids)`. Even if both Blue and Red teams ban the same champion $c$ in blind-ban phase, it counts as **1 match where $c$ was banned** so $\text{BR}_c \in [0\%, 100\%]$ and $\text{PR}_c + \text{BR}_c \le 100\%$.
   - **Denominator:** Total valid matches $M$ in the bucket.

4. **Role Distribution ($\text{RoleShare}_{c,p}$):**
   $$\text{RoleShare}_{c,p} = \frac{\text{Picks}(c, p)}{\sum_{p' \in \{\text{TOP, JGL, MID, BOT, SUP}\}} \text{Picks}(c, p')} \times 100\%$$
