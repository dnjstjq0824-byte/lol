# User Account & Admin Operations Console (9 User Views + 18 Admin Views)
**Document ID:** `BN-RIFTSTATS-UX-PG-USR-01..09, PG-ADM-01..18`  
**Primary Route:** `/account/*, /admin/*`  
**Bound OpenAPI Endpoints:** `GET/POST /api/v1/auth/*, /api/v1/favorites, /api/v1/admin/*`  

---
### 9 User Views & 18 Admin Console Screens
- **User Views (9):** 1. Registration, 2. Login + MFA Challenge, 3. Profile Overview, 4. Favorites Manager, 5. Recent Searches Sync, 6. Riot Account RSO Connection (`APPROVAL_REQUIRED`), 7. Preferences (Locale/Theme/Default Region), 8. Notification Center, 9. RiftStats Pro Subscription Management.
- **Admin Views (18):** 1. Operations Executive Dashboard, 2. Riot API Live Status, 3. Regional Rate-Limit Token Bucket Monitor, 4. API Error & 429 Inspector, 5. Ingestion Job Queue & Dead-Letter Manager, 6. Match Data Explorer, 7. Statistical Aggregation Job Controller, 8. Champion Static Data & Localization Inspector, 9. Patch & Data Dragon Sync Manager, 10. Game Mode & Queue Configurator, 11. Regional Routing & Realm Config, 12. Translation Override Manager, 13. User & RBAC Directory, 14. CLS-Safe Ad Placement Manager, 15. Pro Player Tag Registry, 16. User Report & Abuse Moderation Queue, 17. System Feature Flags (`rso_enabled`, etc.), 18. Immutable Security & Admin Audit Log Viewer.

## State Management, Loading, Empty, Error & Responsive Verification
- **Loading State:** Fixed-dimension CSS grid skeletons prevent layout shift (`CLS < 0.05`).
- **Empty State:** Explicit guidance when a summoner has 0 matches in the selected queue/season or when a filter returns 0 rows.
- **Error State:** Displays structured Problem Details message with one-click `Retry` button; preserves previously loaded DOM data on transient refresh failures.
- **Acceptance Criteria:** Passes Playwright visual layout, DOM click interaction, keyboard navigation, and zero horizontal overflow checks across all 11 mandatory viewports (`320px`, `360px`, `375px`, `390px`, `430px`, `768px`, `1024px`, `1280px`, `1366px`, `1440px`, `1920px`).
