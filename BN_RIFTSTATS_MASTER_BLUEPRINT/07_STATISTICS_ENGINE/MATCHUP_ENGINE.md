# Lane Counter & Team Duo Synergy Statistical Engine
**Document ID:** `BN-RIFTSTATS-STAT-004`  

---

## 1. Direct Lane Matchup Win Rate ($\text{WR}_{A \text{ vs } B, p}$)
Given Champion $A$ in position $p$ and opposing Champion $B$ in the same position $p$ (`team_id_A != team_id_B` and `position_A == position_B` and `position_confidence >= 0.85`):
$$\text{WR}(A \text{ vs } B, p) = \frac{\text{Wins}_A(A_p \text{ vs } B_p)}{\text{Matches}(A_p \text{ vs } B_p)} \times 100\%$$
Additionally computes mean `@15m` lane differentials from `match_timeline_events`:
- $\overline{\text{GD}@15} = \frac{1}{K} \sum_{k=1}^{K} (\text{Gold}_{A,15m,k} - \text{Gold}_{B,15m,k})$
- $\overline{\text{XPD}@15} = \frac{1}{K} \sum_{k=1}^{K} (\text{XP}_{A,15m,k} - \text{XP}_{B,15m,k})$
- $\overline{\text{CSD}@15} = \frac{1}{K} \sum_{k=1}^{K} (\text{CS}_{A,15m,k} - \text{CS}_{B,15m,k})$

## 2. Team Duo Synergy & Normalized Synergy Delta
For Champion $A$ in role $p_1$ and allied Champion $B$ in role $p_2$ ($p_1 \neq p_2$, `team_id_A == team_id_B`):
$$\text{WR}(A_{p_1} + B_{p_2}) = \frac{\text{Wins}(A_{p_1} \land B_{p_2})}{\text{Matches}(A_{p_1} \land B_{p_2})}$$
$$\Delta_{\text{Synergy}}(A_{p_1}, B_{p_2}) = \text{WR}(A_{p_1} + B_{p_2}) - \left(\text{WR}(A_{p_1}) + \text{WR}(B_{p_2}) - 0.5000\right)$$
