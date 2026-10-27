# Summoner Search, Autocomplete & Error States (Public Views 3-5)
**Document ID:** `BN-RIFTSTATS-UX-PG-PUB-03..05`  
**Primary Route:** `/summoners/search`  
**Bound OpenAPI Endpoints:** `GET /api/v1/search/autocomplete, GET /api/v1/accounts/by-riot-id/{region}/{gameName}/{tagLine}`  

---
### Search & Disambiguation Flow
- Parses `gameName#tagLine` (supports whitespace, CJK, Cyrillic, Latin accents).
- When `#tagLine` is omitted, queries OpenSearch for matching `gameName` across the selected region + displays exact regional default (`#KR1`, `#NA1`, `#EUW`, etc.) alongside cross-region matches.
- **Error / Not Found State (View 5):** When `ACCOUNT-V1` returns `404`, renders a clear diagnostic card explaining possible causes (wrong `#tagLine`, wrong region, or uninitialized LoL summoner profile) with one-click buttons to search the same Riot ID across `KR`, `NA`, `EUW`, `VN`, and `TW`.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
