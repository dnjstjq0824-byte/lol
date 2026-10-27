# Champion Builds, Runes, Items, Skills, Counters & Synergies (Public Views 14-19)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-14..19`  
**Primary Route:** `/champions/{championSlug}/build/{role}`  
**Bound OpenAPI Endpoints:** `GET /api/v1/champions/{championKey}/builds, /runes, /items, /counters, /synergies`  

---
### Multi-Choice Build & Matchup Views
- **Builds (View 14):** Presents Top 2 Rune pages, Top 2 Summoner Spell combinations, Skill Priority (`W > Q > E`) + 1-15 grid, Starter Item sets, Boots, Top 5 Core 3-Item builds, and Depth 4 / Depth 5 / Depth 6 situational items.
- **Runes (View 15) & Items (View 16) & Skills (View 17):** Dedicated deep-dive tables sortable by Pick Rate, Win Rate, and Sample Size (`N`), with `LOW_SAMPLE` badges when `10 <= N < 100`.
- **Counters (View 18) & Synergies (View 19):** Displays direct lane matchup win rate, `GD@15`, `XPD@15`, `CSD@15`, Lane Kill Rate, and Allied Role Duo Synergy Delta.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
