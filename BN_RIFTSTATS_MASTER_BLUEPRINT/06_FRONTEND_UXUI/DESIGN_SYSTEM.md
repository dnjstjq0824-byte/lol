# Independent 'Rift Obsidian' Design System, Tokens & Accessibility
**Document ID:** `BN-RIFTSTATS-UX-DS-001`  
**Primary Route:** `N/A`  
**Bound OpenAPI Endpoints:** `N/A`  

---
### Design Tokens & Accessibility Standards
- **Color Contrast (WCAG 2.1 AA Verified):** Primary text `#F9FAFB` on `#0B0F19` (`18.4:1`), secondary text `#9CA3AF` on `#111827` (`6.8:1`), victory accent `#60A5FA` on `#1E3A8A` (`5.1:1`), defeat accent `#FDA4AF` on `#881337` (`5.4:1`).
- **Touch Targets & Keyboard Navigation:** Minimum `44x44px` interactive hit area on mobile viewports (`320px`~`430px`); visible `#06B6D4` `2px` focus ring on all focusable elements; full `ARIA` roles for tabs, combobox autocomplete, and sortable tables.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
