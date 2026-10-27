# Riot API Key Tiers, Security & Lifecycle Policy
**Document ID:** `BN-RIFTSTATS-API-002`  
**Source:** `https://developer.riotgames.com/docs/portal` & `https://developer.riotgames.com/policies/general`  

---

## 1. Comparison of Riot API Key Types

| Attribute | Development API Key | Personal API Key | Production API Key | Riot Sign On (RSO) Client |
|---|---|---|---|---|
| **Intended Use** | Initial local testing & schema exploration | Personal tools, school projects, prototyping for Production Key review | Public commercial/community web platforms (like BN RiftStats) | Verified player OAuth login & private custom match access |
| **Expiration** | **Expires every 24 hours** (requires manual regeneration) | Does not expire daily (tied to single personal project) | Persistent (managed via Developer Portal Group) | Persistent OAuth `client_id` + `client_secret` / Private Key JWT |
| **Application Rate Limit** | **20 req / 1 sec** and **100 req / 2 min** per region | **20 req / 1 sec** and **100 req / 2 min** per region | Starting at **500 req / 10 sec** and **30,000 req / 10 min** per region | Scoped to authenticated user tokens |
| **Public Commercial Launch Permitted?** | **NO (Strictly Prohibited)** | **NO (Not for large-scale public production)** | **YES (Required for public release)** | **YES (Required for account linking)** |

## 2. Non-Negotiable Key Governance Rules
1. **Single Key per Product:** Riot explicitly prohibits registering multiple applications or pooling multiple Personal/Development keys to bypass rate limits (`https://developer.riotgames.com/docs/faqs`). BN RiftStats enforces a strict **single-key-per-environment** architecture.
2. **Server-Side Vault Storage Only:** `RIOT_API_KEY` is stored in AWS Secrets Manager / HashiCorp Vault and injected into the Go API Gateway container memory. It is **never** committed to Git, logged in plaintext, or exposed to the Nuxt browser bundle.
3. **Pre-Approval Operating Mode (`DEV_PROTOTYPE_MODE`):** Until Riot approves the Production API Key, BN RiftStats operates with `APP_ENV=staging_prototype`, enforcing the `20 req/1s, 100 req/2m` token bucket and restricting access to authorized reviewers and automated verification suites.
