# Riot API Global Platform & Regional Routing Specification
**Document ID:** `BN-RIFTSTATS-API-004`  
**Verified Against:** `https://developer.riotgames.com/docs/lol#_routing-values` (2026-10-28)  

---

## 1. Platform vs. Regional Routing Architecture
Riot Games separates endpoints into two distinct routing topologies:
- **Platform Routing (`{platform}.api.riotgames.com`):** Used for `SUMMONER-V4`, `LEAGUE-V4`, `LEAGUE-EXP-V4`, `CHAMPION-MASTERY-V4`, `CHAMPION-V3`, `SPECTATOR-V5`, `LOL-CHALLENGES-V1`, and `LOL-STATUS-V4`.
- **Regional Routing (`{region}.api.riotgames.com`):** Used for `ACCOUNT-V1` and `MATCH-V5`.

## 2. Canonical 16-Region Routing Configuration Table

| UI Region Code | Riot Platform ID | Platform Host | Regional Cluster (MATCH-V5) | Regional Host | Default TagLine | Data Dragon Realm | Primary Locale |
|---|---|---|---|---|---|---|---|
| `KR` | `KR` | `kr.api.riotgames.com` | `ASIA` | `asia.api.riotgames.com` | `KR1` | `kr` | `ko_KR` |
| `NA` | `NA1` | `na1.api.riotgames.com` | `AMERICAS` | `americas.api.riotgames.com` | `NA1` | `na` | `en_US` |
| `EUW` | `EUW1` | `euw1.api.riotgames.com` | `EUROPE` | `europe.api.riotgames.com` | `EUW` | `euw` | `en_GB` |
| `EUNE` | `EUN1` | `eun1.api.riotgames.com` | `EUROPE` | `europe.api.riotgames.com` | `EUNE` | `eune` | `en_GB` |
| `JP` | `JP1` | `jp1.api.riotgames.com` | `ASIA` | `asia.api.riotgames.com` | `JP1` | `jp` | `ja_JP` |
| `BR` | `BR1` | `br1.api.riotgames.com` | `AMERICAS` | `americas.api.riotgames.com` | `BR1` | `br` | `pt_BR` |
| `LAN` | `LA1` | `la1.api.riotgames.com` | `AMERICAS` | `americas.api.riotgames.com` | `LAN` | `lan` | `es_MX` |
| `LAS` | `LA2` | `la2.api.riotgames.com` | `AMERICAS` | `americas.api.riotgames.com` | `LAS` | `las` | `es_AR` |
| `OCE` | `OC1` | `oc1.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `OCE` | `oce` | `en_AU` |
| `TR` | `TR1` | `tr1.api.riotgames.com` | `EUROPE` | `europe.api.riotgames.com` | `TR1` | `tr` | `tr_TR` |
| `RU` | `RU` | `ru.api.riotgames.com` | `EUROPE` | `europe.api.riotgames.com` | `RU1` | `ru` | `ru_RU` |
| `ME` | `ME1` | `me1.api.riotgames.com` | `EUROPE` | `europe.api.riotgames.com` | `ME1` | `me` | `ar_AE` |
| `SG` / `SEA` | `SG2` | `sg2.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `SG2` | `sg` | `en_SG` |
| `PH` | `PH2` | `ph2.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `PH2` | `ph` | `en_PH` |
| `TH` | `TH2` | `th2.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `TH2` | `th` | `th_TH` |
| `TW` | `TW2` | `tw2.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `TW2` | `tw` | `zh_TW` |
| `VN` | `VN2` | `vn2.api.riotgames.com` | `SEA` | `sea.api.riotgames.com` | `VN2` | `vn` | `vi_VN` |

## 3. Strict Routing Validation Guardrails
1. **No Cross-Cluster Match Requests:** Requesting a Korean match (`KR_7123456789`) against `americas.api.riotgames.com` will return `404`. The Go router inspects the Match ID prefix (`KR_`, `NA1_`, `EUW1_`, etc.) or the account's `platform_id` and routes strictly to the mapped `Regional Host`.
2. **Account Lookup Cluster Selection:** For `ACCOUNT-V1`, `americas`, `asia`, and `europe` clusters serve the global account graph. BN RiftStats routes `ACCOUNT-V1` requests to the regional cluster matching the user's selected region (with automatic failover to the nearest secondary cluster on `502/503/504`).
