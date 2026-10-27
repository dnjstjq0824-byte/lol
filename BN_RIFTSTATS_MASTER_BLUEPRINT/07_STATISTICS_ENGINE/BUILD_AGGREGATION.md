# Rune, Item, Skill Order & Summoner Spell Aggregation Specification
**Document ID:** `BN-RIFTSTATS-STAT-003`  

---

## 1. Item & Rune Pick Rate and Win Rate Formulas
For a given champion $c$ and position $p$ with total valid role appearances $N_{c,p} = \text{Picks}(c,p)$:
- **Rune Setup Pick Rate:** $\text{PR}_{\text{rune}} = \frac{\text{GamesWithRuneSignature}(c,p,r)}{N_{c,p}} \times 100\%$
- **Rune Setup Win Rate:** $\text{WR}_{\text{rune}} = \frac{\text{WinsWithRuneSignature}(c,p,r)}{\text{GamesWithRuneSignature}(c,p,r)} \times 100\%$
- **Core 3-Item Build Pick Rate:** $\text{PR}_{\text{core}} = \frac{\text{GamesWithCore3}(c,p,i_1,i_2,i_3)}{\text{TotalGamesReaching3CoreItems}(c,p)} \times 100\%$
- **Depth $k \in \{4,5,6\}$ Item Pick Rate:** Conditional on reaching slot depth $k$, $\text{PR}_{\text{depth } k}(i) = \frac{\text{GamesPurchasingAtDepth}(c,p,k,i)}{\text{TotalGamesReachingDepth}(c,p,k)} \times 100\%$.
- **Riot Arena Policy Guard:** For `queue_id IN (1700, 1710)` (Arena), item and augment win-rate aggregation is strictly disabled (`FORBIDDEN_BY_RIOT_POLICY`).
