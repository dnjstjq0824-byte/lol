# Riot Data Dragon & Static Game Asset Pipeline
**Document ID:** `BN-RIFTSTATS-API-007`  
**Verified Against:** `https://developer.riotgames.com/docs/lol#_data-dragon` & `https://ddragon.leagueoflegends.com/api/versions.json` (Latest Live Version: `16.21.1`)  

---

## 1. Official Data Dragon Endpoints & Version Synchronization
1. **Global Version List:** `https://ddragon.leagueoflegends.com/api/versions.json` (Polled every `30 minutes` by Go Static Data Sync Worker).
2. **Regional Realm Version File:** `https://ddragon.leagueoflegends.com/realms/{region}.json` (e.g., `kr.json`, `na.json`, `euw.json`) — verifies which Data Dragon version is active in each region.
3. **Supported Locales:** `https://ddragon.leagueoflegends.com/cdn/languages.json` (`ko_KR`, `en_US`, `ja_JP`, `zh_CN`, `zh_TW`, `de_DE`, `es_ES`, `fr_FR`, `pt_BR`, `ru_RU`, `tr_TR`, `vi_VN`, `th_TH`, `pl_PL`, etc.).
4. **Core JSON Manifests (Version `16.21.1`):**
   - Champions Summary: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/champion.json`
   - Champion Full Detail: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/champion/{ChampionKey}.json`
   - Mode-Specific Champions (e.g. League Classic): `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/mode/classic/champion.json`
   - Items: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/item.json`
   - Runes Reforged: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/runesReforged.json`
   - Summoner Spells: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/summoner.json`
   - Profile Icons: `https://ddragon.leagueoflegends.com/cdn/16.21.1/data/{locale}/profileicon.json`
5. **Auxiliary Official Game Constants:**
   - Queues: `https://static.developer.riotgames.com/docs/lol/queues.json`
   - Maps: `https://static.developer.riotgames.com/docs/lol/maps.json`
   - Game Modes: `https://static.developer.riotgames.com/docs/lol/gameModes.json`

## 2. Handling Patch-Day Lag & Missing Asset Resilience
As documented by Riot (*"Updating Data Dragon after each League of Legends patch is a manual process, so it is not always updated immediately after a patch"*):
1. **Patch Normalization (`major.minor`):** A `MATCH-V5` payload may report `gameVersion: "16.22.612.3456"` two hours before `16.22.1` is published to `versions.json`.
2. **Fallback Resolution:** When a match's `major.minor` (`16.22`) is newer than the latest Data Dragon version (`16.21.1`), the asset resolver maps champion/item/spell IDs against the most recent available Data Dragon version (`16.21.1`) and logs an `UNRESOLVED_ASSET_ID` warning in `source_freshness` if a brand-new item/champion ID is encountered, rendering a clean fallback icon (`unknown_item.svg`) rather than a broken `404` image.
3. **Zero Third-Party Asset Theft:** All game icons, splashes, and ability images are served from Data Dragon or mirrored onto BN RiftStats's own S3/CloudFront CDN directly from official Data Dragon tarballs (`dragontail-16.21.1.tgz`). OP.GG CDN URLs (`opgg-static.akamaized.net`) are strictly prohibited.
