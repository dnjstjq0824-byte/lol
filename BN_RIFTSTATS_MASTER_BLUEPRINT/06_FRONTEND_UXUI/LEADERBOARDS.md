# Global & Regional Leaderboards (Public View 20)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-20`  
**Primary Route:** `/leaderboards/{region}`  
**Bound OpenAPI Endpoints:** `GET /api/v1/leaderboards/{region}`  

---
### Leaderboard Modes
- **Ranked Solo/Duo & Ranked Flex Ladders:** Official Riot `LEAGUE-V4` / `LEAGUE-EXP-V4` standings with Rank `#`, Riot ID (`gameName#tagLine`), Verified Pro badge, Tier Emblem, LP, Level, Most Played 3 Champions, and Win Rate bar.
- **Champion Mastery & Summoner Level Rankings:** Filtered by region and champion, displaying clear data-source attribution (`Indexed Public Summoners`).

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
