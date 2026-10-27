# Match History Cards & Filter Engine (Public View 9)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-09`  
**Primary Route:** `/summoners/{region}/{gameName}-{tagLine}/matches`  
**Bound OpenAPI Endpoints:** `GET /api/v1/summoners/{region}/by-puuid/{puuid}/matches`  

---
### Match Card Responsive & Data Rules
- Displays Queue Name, Relative Timestamp, Result (`Victory` blue, `Defeat` rose, `Remake` neutral slate), Duration, Champion Avatar + Level, Summoner Spells, Primary Keystone + Secondary Rune Tree, K/D/A, KDA Ratio, Kill Participation (`P/Kill %`), CS (`Total (CS/m)`), Average Match Tier, 6 Item Slots + Trinket + Season 2026 Role Quest Item, Multi-Kill Badge, and Rift Score Rank (`MVP`, `ACE`, `2nd`~`10th`).
- Prevents icon overlap and text truncation across all 11 viewports (`320px`~`1920px`).

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
