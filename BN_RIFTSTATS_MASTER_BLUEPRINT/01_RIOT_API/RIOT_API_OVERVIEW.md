# Riot Games Official API — Architecture & Integration Overview
**Document ID:** `BN-RIFTSTATS-API-001`  
**Verified Against:** `https://developer.riotgames.com/apis` & `https://developer.riotgames.com/docs/lol` (2026-10-28)  

---

## 1. Mandatory Riot API Inventory (10 Official API Families)
BN RiftStats integrates exclusively with documented, official Riot Games APIs over HTTPS (`TLS 1.3`). All requests are executed server-side by the Go Riot API Gateway (`internal/riotclient`) using the `X-Riot-Token` header (or `Authorization: Bearer {accessToken}` for RSO-authenticated endpoints).

| # | API Family | Routing Type | Primary Purpose in BN RiftStats | Official Support Status |
|---|---|---|---|---|
| 1 | `ACCOUNT-V1` | **Regional** (`americas`, `asia`, `europe`, `sea`) | Resolve `gameName#tagLine` ↔ `puuid`, and fetch authenticated user via RSO `/accounts/me` | **ACTIVE (Authoritative Identity API)** |
| 2 | `SUMMONER-V4` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch `profileIconId`, `summonerLevel`, `revisionDate`, and `id` (encryptedSummonerId) by `puuid` | **ACTIVE (`/by-name` deprecated; use `/by-puuid`)** |
| 3 | `MATCH-V5` | **Regional** (`americas`, `asia`, `europe`, `sea`) | Fetch match ID lists by `puuid`, full match payloads (`/matches/{matchId}`), and minute-by-minute timelines (`/matches/{matchId}/timeline`) | **ACTIVE (Core Match & Analytics Source)** |
| 4 | `LEAGUE-V4` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch Ranked Solo/Duo and Ranked Flex entries by `puuid` or `summonerId`, plus Challenger, Grandmaster, and Master league ladders | **ACTIVE (Official Ranked Ladder Source)** |
| 5 | `LEAGUE-EXP-V4` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Paginated enumeration of league entries across all tiers (`CHALLENGER` down to `IRON`) and divisions (`I`~`IV`) for leaderboard & seed crawling | **ACTIVE (Experimental/High-Volume Ladder Endpoint)** |
| 6 | `CHAMPION-MASTERY-V4` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch player's champion mastery scores, levels, points, and top mastery champions by `puuid` | **ACTIVE (PUUID-native)** |
| 7 | `CHAMPION-V3` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch weekly free-to-play champion rotation (`freeChampionIds`, `freeChampionIdsForNewPlayers`) | **ACTIVE** |
| 8 | `SPECTATOR-V5` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch active live game info (`/active-games/by-summoner/{encryptedPUUID}`) and featured games | **ACTIVE (Replaced v4; uses PUUID)** |
| 9 | `LOL-CHALLENGES-V1` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch player challenge progression, titles, crystal level, and challenge leaderboards by `puuid` | **ACTIVE** |
| 10 | `LOL-STATUS-V4` | **Platform** (`kr`, `na1`, `euw1`, etc.) | Fetch live shard status, maintenance schedules, and active incidents per platform region | **ACTIVE** |

---

## 2. Detailed Per-API Technical Specifications

### 2.1 `ACCOUNT-V1`
- **Purpose:** Authoritative resolution between Riot ID (`gameName` + `tagLine`) and globally unique `puuid` (78-character encrypted string).
- **Endpoints:**
  1. `GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}`
  2. `GET /riot/account/v1/accounts/by-puuid/{puuid}`
  3. `GET /riot/account/v1/accounts/me` (Requires RSO Bearer token)
  4. `GET /riot/account/v1/active-shards/by-game/{game}/by-puuid/{puuid}`
- **Routing Host:** Regional (`americas.api.riotgames.com`, `asia.api.riotgames.com`, `europe.api.riotgames.com`, `sea.api.riotgames.com`). Note: `/by-riot-id` and `/by-puuid` return identical account data across `americas`, `asia`, and `europe`; Riot recommends routing to the cluster closest to the backend server.
- **Response Schema (`AccountDto`):**
  ```json
  {
    "puuid": "78-char-encrypted-puuid",
    "gameName": "Hide on bush",
    "tagLine": "KR1"
  }
  ```
- **Cache Policy:** Redis TTL `3600s` for `riotId -> puuid` mapping; PostgreSQL `riot_accounts` table stores persistent `puuid` as Primary Key and records historical name changes in `riot_id_history`.
- **Error Handling:** `404 Not Found` → Return structured `RIOT_ACCOUNT_NOT_FOUND` (suggest checking `#tagLine` or region); `429 Too Many Requests` → Respect `Retry-After` header.

### 2.2 `SUMMONER-V4`
- **Purpose:** Retrieve platform-specific LoL summoner profile metadata (`profileIconId`, `summonerLevel`, `revisionDate`, `id`, `accountId`).
- **Endpoints:**
  1. `GET /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}` (**Primary**)
  2. `GET /lol/summoner/v4/summoners/{encryptedSummonerId}`
  3. `GET /lol/summoner/v4/summoners/me` (RSO Bearer token)
  - *Note:* `GET /lol/summoner/v4/summoners/by-name/{summonerName}` is **DEPRECATED** per Riot's November 2023 Riot ID migration and MUST NOT be called by BN RiftStats.
- **Routing Host:** Platform (`kr.api.riotgames.com`, `na1.api.riotgames.com`, `euw1.api.riotgames.com`, etc.).
- **Response Schema (`SummonerDto`):**
  ```json
  {
    "id": "encrypted-summoner-id",
    "accountId": "encrypted-account-id",
    "puuid": "78-char-encrypted-puuid",
    "profileIconId": 6,
    "revisionDate": 1761600000000,
    "summonerLevel": 946
  }
  ```
- **Cache Policy:** Redis TTL `300s` (5 minutes); bypassed on explicit user `Update` click if `now - last_updated >= 120s`.

### 2.3 `MATCH-V5`
- **Purpose:** Discover recent Match IDs for a player and retrieve full 10-player match statistics and event timelines.
- **Endpoints:**
  1. `GET /lol/match/v5/matches/by-puuid/{puuid}/ids?start={0}&count={20}&queue={queueId}&type={type}&startTime={epochSec}&endTime={epochSec}`
  2. `GET /lol/match/v5/matches/{matchId}`
  3. `GET /lol/match/v5/matches/{matchId}/timeline`
- **Routing Host:** Regional (`americas`, `asia`, `europe`, `sea`).
  - `NA1`, `BR1`, `LA1`, `LA2` → `americas.api.riotgames.com`
  - `KR`, `JP1` → `asia.api.riotgames.com`
  - `EUW1`, `EUN1`, `TR1`, `RU`, `ME1` → `europe.api.riotgames.com`
  - `OC1`, `PH2`, `SG2`, `TH2`, `TW2`, `VN2` → `sea.api.riotgames.com`
- **Cache Policy:** Match details and timelines are immutable once completed (`PostgreSQL` permanent storage + Redis/S3 cache; never re-fetched from Riot once `match_ingestion_status = 'COMPLETED'`).
- **Data Usage Restriction:** Custom games (`queueId = 0` or custom tournament lobbies) MUST NOT be publicly displayed in a player's match history unless the player explicitly authenticated via RSO and opted in.

### 2.4 `LEAGUE-V4` & `LEAGUE-EXP-V4`
- **Purpose:** Retrieve official Ranked Solo/Duo (`RANKED_SOLO_5x5`) and Ranked Flex (`RANKED_FLEX_SR`) tier, division, LP, wins, losses, streaks, and apex tier ladders.
- **Endpoints:**
  1. `GET /lol/league/v4/entries/by-puuid/{encryptedPUUID}` (Supported modern endpoint)
  2. `GET /lol/league/v4/entries/by-summoner/{encryptedSummonerId}`
  3. `GET /lol/league/v4/challengerleagues/by-queue/{queue}`
  4. `GET /lol/league/v4/grandmasterleagues/by-queue/{queue}`
  5. `GET /lol/league/v4/masterleagues/by-queue/{queue}`
  6. `GET /lol/league/v4/entries/{queue}/{tier}/{division}?page={page}`
  7. `GET /lol/league-exp/v4/entries/{queue}/{tier}/{division}?page={page}`
- **Routing Host:** Platform.
- **Cache Policy:** Player ranked entries cached `180s` in Redis; Apex ladders (`Challenger`/`GM`/`Master`) refreshed every `10 minutes` via background cron job.

### 2.5 `CHAMPION-MASTERY-V4`
- **Purpose:** Retrieve player's champion mastery level, points, tokens, and last play timestamp.
- **Endpoints:**
  1. `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}`
  2. `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}/top?count={count}`
  3. `GET /lol/champion-mastery/v4/champion-masteries/by-puuid/{encryptedPUUID}/by-champion/{championId}`
  4. `GET /lol/champion-mastery/v4/scores/by-puuid/{encryptedPUUID}`
- **Routing Host:** Platform.
- **Cache Policy:** Redis TTL `900s` (15 minutes).

### 2.6 `CHAMPION-V3`
- **Purpose:** Weekly free champion rotation (`GET /lol/platform/v3/champion-rotations`).
- **Routing Host:** Platform.
- **Cache Policy:** Redis TTL `3600s` (1 hour).

### 2.7 `SPECTATOR-V5`
- **Purpose:** Real-time active match inspection.
- **Endpoints:**
  1. `GET /lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}`
  2. `GET /lol/spectator/v5/featured-games`
- **Routing Host:** Platform.
- **Cache Policy:** Short Redis TTL `15s` with singleflight deduplication so 50 concurrent viewers of a pro streamer's live game generate only 1 Riot API call per 15 seconds.
- **Error Handling:** Returns `404 Not Found` when the player is not currently in an active match—mapped to HTTP `200 OK` with `{"inGame": false}` in BN RiftStats API so client logs remain clean.

### 2.8 `LOL-CHALLENGES-V1`
- **Purpose:** Retrieve player challenge crystal level, points, and category percentiles (`GET /lol/challenges/v1/player-data/{puuid}`).
- **Routing Host:** Platform.
- **Cache Policy:** Redis TTL `1800s` (30 minutes).

### 2.9 `LOL-STATUS-V4`
- **Purpose:** Retrieve live server maintenance and incident banners (`GET /lol/status/v4/platform-data`).
- **Routing Host:** Platform.
- **Cache Policy:** Polled every `120s` per platform by background worker; served from Redis cache (`TTL 120s`).
