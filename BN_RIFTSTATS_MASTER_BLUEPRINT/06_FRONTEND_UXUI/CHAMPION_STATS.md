# Champion Tier List & Champion Detail Overview (Public Views 12-13, 21-22)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-12..13, 21..22`  
**Primary Route:** `/champions/tier-list`  
**Bound OpenAPI Endpoints:** `GET /api/v1/champions/tier-list, GET /api/v1/modes/{modeSlug}/stats, GET /api/v1/patches`  

---
### Champion Tier Table & Filter Controls
- Filter controls: `Server` (16 regions + `GLOBAL`), `Tier Bracket` (`Emerald+` default, `Challenger` down to `Iron`), `Patch` (`16.21`, `16.20`...), `Game Mode` (`Ranked Solo/Duo`, `Flex`, `ARAM`, `Arena`, `ARAM: Mayhem`, `URF`), and `Role` (`All`, `Top`, `Jungle`, `Mid`, `Bottom`, `Support`).
- Sortable columns (`tabular-nums`): `Rank`, `Champion`, `Tier` (`OP`, `1`~`5`), `Role`, `Win Rate`, `Pick Rate`, `Ban Rate`, `Weak Against (Top 3 Counters)`, and `Sample Size (N)`.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
