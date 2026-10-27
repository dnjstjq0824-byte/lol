# BN RiftStats — STRIDE Threat Model & Mitigation Matrix
**Document ID:** `BN-RIFTSTATS-ARCH-005`  

---

## 1. STRIDE Threat Analysis

| Threat ID | Category | Threat Vector | Impact | Architectural Mitigation | Verification Test |
|---|---|---|---|---|---|
| `TM-01` | **Information Disclosure** | Leakage of `RIOT_API_KEY` into client JS bundle, browser network tab, or error logs | Immediate Riot API key revocation | Key isolated in Go backend memory; stripped from all logs/errors; CI secret scanning | `TEST-SEC-01` |
| `TM-02` | **Denial of Service (Rate Limit Exhaustion)** | Malicious bot triggers thousands of `Update` requests for random summoners to exhaust Riot 429 quota | Global regional outage for legitimate users | Per-IP rate limiter (`30 req/min`) + per-summoner `120s` refresh cooldown + P0/P1/P2 priority shedding | `TEST-SEC-02` |
| `TM-03` | **Spoofing / Account Takeover** | Attacker attempts to claim ownership of a pro player's Riot ID (`Hide on bush#KR1`) without authorization | Impersonation on platform profile | Verified badge requires cryptographic RSO OAuth (`/riot/account/v1/accounts/me`); unverified bookmarks are read-only Favorites | `TEST-SEC-03` |
| `TM-04` | **Tampering / SQL Injection** | Malicious characters in `gameName#tagLine` (e.g., `' OR 1=1--`) or sort parameters | Database compromise | 100% `pgx` parameterized queries + strict enum validation on sort/filter parameters | `TEST-SEC-04` |
| `TM-05` | **Cross-Site Scripting (XSS)** | Crafted Riot `gameName` containing `<script>` or SVG payload rendered in match history | Session theft | Vue 3 automatic text escaping + strict `Content-Security-Policy` (`script-src 'self'`) | `TEST-SEC-05` |
| `TM-06` | **Elevation of Privilege** | Standard user attempts to access `/api/v1/admin/*` endpoints | Unauthorized system config or key inspection | Middleware enforces JWT signature + `role IN ('SUPER_ADMIN', 'OPS_ADMIN')` + MFA verification claim | `TEST-SEC-06` |
