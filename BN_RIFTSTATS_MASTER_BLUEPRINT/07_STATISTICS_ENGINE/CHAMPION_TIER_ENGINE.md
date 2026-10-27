# Champion Tier Scoring Engine (Wilson Shrinkage + Z-Score Composite)
**Document ID:** `BN-RIFTSTATS-STAT-002`  

---

## 1. Preventing Low-Sample Distortion via Wilson Score Lower Bound
To prevent a niche champion with `18 wins / 25 games (72% raw WR)` from outranking a meta staple with `26,000 wins / 50,000 games (52% WR)`, BN RiftStats computes the **95% Wilson Score Confidence Interval Lower Bound** ($z = 1.95996$):
$$\text{WR}_{\text{Wilson}}(c,p) = \frac{\hat{p} + \frac{z^2}{2n} - z \sqrt{\frac{\hat{p}(1 - \hat{p}) + \frac{z^2}{4n}}{n}}}{1 + \frac{z^2}{n}}$$
where $\hat{p} = \text{Wins}(c,p) / n$ and $n = \text{Picks}(c,p)$.

## 2. Composite Tier Score & Grade Cutoffs
Within each `(patch, region, tier_bracket, queue_id, position)` slice (filtered to champions with $\text{RoleShare}_{c,p} \ge 5\%$ and $n \ge 100$):
$$S_{c,p} = 0.55 \cdot Z\big(\text{WR}_{\text{Wilson}}(c,p)\big) + 0.30 \cdot Z\big(\ln(1 + \text{PR}_{c,p})\big) + 0.15 \cdot Z\big(\text{BR}_c\big) \cdot \mathbb{I}(\hat{p} \ge 0.500)$$
- **`OP` Tier:** Top $4\%$ of role ($S_{c,p} \ge +1.75$ and $\text{PR}_{c,p} \ge 3.0\%$)
- **`Tier 1`:** $S_{c,p} \in [+0.85, +1.75)$
- **`Tier 2`:** $S_{c,p} \in [+0.15, +0.85)$
- **`Tier 3`:** $S_{c,p} \in [-0.55, +0.15)$
- **`Tier 4`:** $S_{c,p} \in [-1.25, -0.55)$
- **`Tier 5`:** $S_{c,p} < -1.25$
- *Policy Note:* This Tier Grade is an analytical champion meta classification and NEVER evaluates or replaces a player's official Riot Ranked Ladder standing.
