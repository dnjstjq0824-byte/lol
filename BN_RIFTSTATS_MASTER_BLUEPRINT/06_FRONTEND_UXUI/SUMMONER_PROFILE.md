# Summoner Profile, Ranked Solo/Flex & Mastery (Public Views 6-8, 23)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-06..08, 23`  
**Primary Route:** `/summoners/{region}/{gameName}-{tagLine}`  
**Bound OpenAPI Endpoints:** `GET /api/v1/summoners/{region}/by-puuid/{puuid}, /ranked, /mastery, /spectator`  

---
### ASCII Wireframe (Desktop 1440px 2-Column Layout)
```text
+-----------------------------------------------------------------------------------+
| [Icon Lv.946] Hide on bush #KR1  [Pro: T1]  Ladder Rank #83 (Top 0.0026%)         |
| [S2025 Master] [S2024 GM]   [Update (Cooldown 120s)] [Tier Graph] [☆ Favorite]    |
| Last updated: 3 minutes ago                                                       |
+-----------------------------------------------------------------------------------+
| [Summary] | [Champions] | [Mastery] | [Live Game (Green Pulse if Active)]         |
+------------------------+----------------------------------------------------------+
| LEFT RAIL (320px)      | MAIN MATCH STREAM (Fluid / 740px+)                       |
| - Ranked Solo Card     | - Queue Filter: [All][Solo][Flex][ARAM][Arena][Modes v]  |
|   Challenger 2,256 LP  | - 20-Game Summary: 20G 13W 7L (65%) | 2.36:1 KDA | 46%KP |
|   448W 371L (55%)      | - Match Card #1 (Victory | 16m 51s | Galio 7/4/6 | 4th)  |
| - Ranked Flex Card     | - Match Card #2 (Defeat  | 24m 12s | LeBlanc 2/6/2)      |
| - S2026 Playstyle      | - Match Card #3 (Victory | 24m 39s | Rumble 8/4/9 | 2nd) |
| - Most Played Champs   | - [Show More Matches (Cursor Pagination)]                |
| - Recently Played With |                                                          |
+------------------------+----------------------------------------------------------+
```

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
