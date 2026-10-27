# Riot Sign On (RSO) OAuth 2.0 / OIDC Integration Specification
**Document ID:** `BN-RIFTSTATS-API-006`  
**Verified Against:** `https://developer.riotgames.com/docs/lol#_rso-integration`  

---

## 1. Approval Status & Dual-Mode Design
- **Policy Requirement:** RSO access requires an approved **Production API Key** and a dedicated RSO OAuth Client (`client_id`) issued by Riot Games Developer Relations.
- **Phase Gating (`FEATURE_FLAG_RSO_ENABLED`):**
  - **When `false` (Pre-RSO Approval):** Users register/login via BN RiftStats native email/password or Google/Discord OAuth and bookmark any public Riot ID (`gameName#tagLine`). No Riot credentials are ever requested.
  - **When `true` (Post-RSO Approval):** Users click **"Connect Riot Account (Official RSO)"** to verify ownership of their Riot account and optionally opt in to private custom match history indexing.

## 2. Official RSO OAuth 2.0 Authorization Code Flow
1. **Authorization Redirect:**
   ```text
   https://auth.riotgames.com/authorize?client_id={CLIENT_ID}&redirect_uri={HTTPS_CALLBACK_URL}&response_type=code&scope=openid+offline_access&state={CSRF_TOKEN}
   ```
2. **Token Exchange (`POST https://auth.riotgames.com/token`):**
   - Authenticated via `client_secret_basic` or `private_key_jwt`.
   - Exchanges `code` for `access_token`, `refresh_token`, and `id_token`.
3. **Identity Verification (`GET /riot/account/v1/accounts/me`):**
   - Called with header `Authorization: Bearer {access_token}` against nearest regional cluster (`asia.api.riotgames.com`, `americas.api.riotgames.com`, or `europe.api.riotgames.com`).
   - Returns verified `puuid`, `gameName`, and `tagLine`, which is cryptographically bound to `auth_identities` (`provider = 'riot_rso'`).
4. **Zero Password Collection:** BN RiftStats **never** collects, proxies, or stores Riot account passwords under any circumstance.
