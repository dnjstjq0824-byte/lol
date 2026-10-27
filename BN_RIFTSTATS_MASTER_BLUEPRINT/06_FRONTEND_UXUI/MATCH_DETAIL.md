# Match Detail Accordion & Participant Comparison (Public Views 10-11)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-10..11`  
**Primary Route:** `/matches/{region}/{matchId}`  
**Bound OpenAPI Endpoints:** `GET /api/v1/matches/{region}/{matchId}, GET /api/v1/matches/{region}/{matchId}/timeline`  

---
### 4-Tab Expanded Match Detail Specification
1. **Overview Tab:** Full 10-player scoreboard (5 Blue Team, 5 Red Team) showing Champion, Spells, Runes, Riot ID, Tier, Rift Score, KDA, Damage Dealt to Champions (with relative horizontal bar), Damage Taken, Control Wards, Wards Placed/Killed, CS, Gold, and Items, plus Team Objective summary bar (Baron, Dragon, Atakhan, Voidgrubs, Herald, Towers).
2. **Team Analysis Tab:** Head-to-head team distribution bars for Kills, Gold Earned, Damage Dealt, Damage Taken, Creep Score, and Vision Score.
3. **Build Tab:** Personal chronological item purchase timeline grouped by minute (`0m`, `5m`, `11m`...), 1-18 skill level-up matrix, and full selected rune tree with stat shards.
4. **Timeline / Etc Tab:** Interactive SVG Gold Advantage Curve, XP Advantage Curve, and @15m Lane Matchup comparison (`GD@15`, `XPD@15`, `CSD@15`). Gracefully shows a notice if `MATCH-V5 Timeline` is not available for that game mode.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
