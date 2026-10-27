# Global Home & Server Selection (Public Views 1-2)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-01..02`  
**Primary Route:** `/`  
**Bound OpenAPI Endpoints:** `GET /api/v1/regions, GET /api/v1/champions/tier-list?limit=5, GET /api/v1/patches`  

---
### ASCII Wireframe (Desktop 1440px & Mobile 375px)
```text
+-----------------------------------------------------------------------------------+
| [BN RiftStats Logo]  [Home] [Champions] [Game Modes] [Leaderboards] [Multi-Search]|
|                                          [Server: KR v] [ko_KR v] [Dark] [Sign In]|
+-----------------------------------------------------------------------------------+
|                           BN RIFTSTATS — PATCH 16.21                              |
|         +---------------+-----------------------------------------------+         |
|         | Region: KR v  | Search Riot ID (e.g. Hide on bush#KR1)   [Go] |         |
|         +---------------+-----------------------------------------------+         |
|         [Recent Searches / Favorites / Instant Champion Autocomplete Popover]     |
+-----------------------------------------------------------------------------------+
|  [Patch 16.21 Meta Tier Summary (Top / Jungle / Mid / Bot / Support)]             |
|  [Live Server Status & Data Freshness Banner]      [CLS-Safe Ad Slot 728x90]      |
+-----------------------------------------------------------------------------------+
```
### Component Tree
- `<AppLayout>` → `<GlobalHeader>` (`<RegionSelector />`, `<LocaleSelector />`, `<ThemeToggle />`, `<UserMenu />`)
- `<HomeHeroSearch>` (`<RiotIdSearchBar />`, `<SearchAutocompletePopover />`)
- `<HomePatchMetaGrid>` (`<RoleTierMiniCard />`, `<PatchHighlightList />`)
- `<AdSlotContainer placementId="home-leaderboard" />`

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
